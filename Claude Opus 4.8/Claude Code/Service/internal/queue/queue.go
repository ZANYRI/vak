package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const (
	streamMain  = "billing:jobs"
	streamDead  = "billing:jobs:dead"
	groupName   = "workers"
)

// ErrDuplicate indicates a job with the same idempotency key already exists.
var ErrDuplicate = errors.New("duplicate job")

// Queue is a Redis Streams broker backed by a Postgres jobs table for
// lifecycle/observability. The API publishes; workers consume.
type Queue struct {
	rdb            *redis.Client
	db             *pgxpool.Pool
	defaultMaxTry  int
}

// New connects to Redis using a redis:// URL and ensures the consumer group.
func New(ctx context.Context, redisURL string, db *pgxpool.Pool, defaultMaxAttempts int) (*Queue, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}
	rdb := redis.NewClient(opt)
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}
	q := &Queue{rdb: rdb, db: db, defaultMaxTry: defaultMaxAttempts}
	if err := q.ensureGroup(ctx); err != nil {
		return nil, err
	}
	return q, nil
}

func (q *Queue) ensureGroup(ctx context.Context) error {
	// MKSTREAM creates the stream if absent; ignore BUSYGROUP on re-create.
	err := q.rdb.XGroupCreateMkStream(ctx, streamMain, groupName, "0").Err()
	if err != nil && !isBusyGroup(err) {
		return fmt.Errorf("create consumer group: %w", err)
	}
	return nil
}

func isBusyGroup(err error) bool {
	return err != nil && err.Error() == "BUSYGROUP Consumer Group name already exists"
}

// Close releases the Redis connection.
func (q *Queue) Close() error { return q.rdb.Close() }

// Ping verifies Redis connectivity (used by readiness checks).
func (q *Queue) Ping(ctx context.Context) error { return q.rdb.Ping(ctx).Err() }

// Publish records a job in Postgres and enqueues it on the Redis stream.
// If opts.IdempotencyKey is set and already used, it returns ErrDuplicate.
func (q *Queue) Publish(ctx context.Context, jobType string, payload any, opts PublishOptions) (string, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}
	maxTry := opts.MaxAttempts
	if maxTry <= 0 {
		maxTry = q.defaultMaxTry
	}

	var jobID string
	var idem *string
	if opts.IdempotencyKey != "" {
		idem = &opts.IdempotencyKey
	}

	err = q.db.QueryRow(ctx, `
		INSERT INTO jobs (type, payload, status, max_attempts, idempotency_key)
		VALUES ($1, $2, 'queued', $3, $4)
		RETURNING id`, jobType, raw, maxTry, idem).Scan(&jobID)
	if err != nil {
		if isUniqueViolation(err) {
			return "", ErrDuplicate
		}
		return "", fmt.Errorf("insert job: %w", err)
	}

	job := Job{ID: jobID, Type: jobType, Payload: raw, MaxAttempts: maxTry}
	if opts.IdempotencyKey != "" {
		job.IdempotencyKey = opts.IdempotencyKey
	}
	if err := q.xadd(ctx, streamMain, job); err != nil {
		return "", err
	}
	return jobID, nil
}

func (q *Queue) xadd(ctx context.Context, stream string, job Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return q.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: map[string]any{"job": string(data)},
	}).Err()
}

// QueueDepth returns the number of pending (un-acked) entries on the main stream.
func (q *Queue) QueueDepth(ctx context.Context) (int64, error) {
	return q.rdb.XLen(ctx, streamMain).Result()
}

func isUniqueViolation(err error) bool {
	var pgErr interface{ SQLState() string }
	if errors.As(err, &pgErr) {
		return pgErr.SQLState() == "23505"
	}
	return false
}

// markStatus updates the jobs table lifecycle column.
func (q *Queue) markStatus(ctx context.Context, jobID, status, lastErr string, attempts int) {
	_, _ = q.db.Exec(ctx, `
		UPDATE jobs SET status = $2, attempts = $3, last_error = $4, updated_at = now()
		WHERE id = $1`, jobID, status, attempts, lastErr)
}

// jobExistsCompleted reports whether the job row is already completed (idempotency).
func (q *Queue) jobExistsCompleted(ctx context.Context, jobID string) (bool, error) {
	var status string
	err := q.db.QueryRow(ctx, `SELECT status FROM jobs WHERE id = $1`, jobID).Scan(&status)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return status == "completed", nil
}

// backoff computes exponential backoff for a given attempt number.
func backoff(attempt int) time.Duration {
	d := time.Duration(1<<uint(attempt)) * time.Second
	if d > 5*time.Minute {
		d = 5 * time.Minute
	}
	return d
}
