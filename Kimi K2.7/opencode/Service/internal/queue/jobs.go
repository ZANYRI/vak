package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"billing-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// JobStore persists job state in PostgreSQL.
type JobStore struct {
	pool *pgxpool.Pool
}

// NewJobStore creates a job store.
func NewJobStore(pool *pgxpool.Pool) *JobStore {
	return &JobStore{pool: pool}
}

// CreateJob inserts a queued job and returns its ID.
func (js *JobStore) CreateJob(ctx context.Context, queue string, payload map[string]interface{}, maxAttempts int) (uuid.UUID, error) {
	id := uuid.New()
	if maxAttempts <= 0 {
		maxAttempts = 3
	}
	_, err := js.pool.Exec(ctx,
		`INSERT INTO jobs (id, queue, payload, status, attempts, max_attempts, run_after)
		 VALUES ($1, $2, $3, 'queued', 0, $4, NOW())`,
		id, queue, payload, maxAttempts)
	if err != nil {
		return uuid.Nil, fmt.Errorf("insert job: %w", err)
	}
	return id, nil
}

// GetJob returns a job by id.
func (js *JobStore) GetJob(ctx context.Context, id uuid.UUID) (*models.Job, error) {
	var j models.Job
	var payload []byte
	err := js.pool.QueryRow(ctx,
		`SELECT id, queue, payload, status, attempts, max_attempts, run_after, error_message, dead_at, created_at, updated_at
		 FROM jobs WHERE id=$1`, id).Scan(
		&j.ID, &j.Queue, &payload, &j.Status, &j.Attempts, &j.MaxAttempts, &j.RunAfter, &j.ErrorMessage, &j.DeadAt, &j.CreatedAt, &j.UpdatedAt)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(payload, &j.Payload)
	return &j, nil
}

// MarkRunning sets status running and increments attempts.
func (js *JobStore) MarkRunning(ctx context.Context, id uuid.UUID) error {
	_, err := js.pool.Exec(ctx,
		`UPDATE jobs SET status='running', attempts=attempts+1, updated_at=NOW() WHERE id=$1`, id)
	return err
}

// MarkCompleted sets status completed.
func (js *JobStore) MarkCompleted(ctx context.Context, id uuid.UUID) error {
	_, err := js.pool.Exec(ctx,
		`UPDATE jobs SET status='completed', updated_at=NOW() WHERE id=$1`, id)
	return err
}

// MarkFailed handles failure, retry with exponential backoff, or dead.
func (js *JobStore) MarkFailed(ctx context.Context, id uuid.UUID, errMsg string) error {
	job, err := js.GetJob(ctx, id)
	if err != nil {
		return err
	}
	if job.Attempts >= job.MaxAttempts {
		deadAt := time.Now()
		_, err = js.pool.Exec(ctx,
			`UPDATE jobs SET status='dead', error_message=$1, dead_at=$2, updated_at=NOW() WHERE id=$3`,
			errMsg, deadAt, id)
		return err
	}
	backoff := time.Duration(math.Pow(2, float64(job.Attempts))) * time.Second
	if backoff > 5*time.Minute {
		backoff = 5 * time.Minute
	}
	_, err = js.pool.Exec(ctx,
		`UPDATE jobs SET status='retrying', error_message=$1, run_after=$2, updated_at=NOW() WHERE id=$3`,
		errMsg, time.Now().Add(backoff), id)
	return err
}

// MarkDead forces a job to dead status.
func (js *JobStore) MarkDead(ctx context.Context, id uuid.UUID, errMsg string) error {
	_, err := js.pool.Exec(ctx,
		`UPDATE jobs SET status='dead', error_message=$1, dead_at=NOW(), updated_at=NOW() WHERE id=$2`,
		errMsg, id)
	return err
}

// ListRetryable returns jobs ready to be retried.
func (js *JobStore) ListRetryable(ctx context.Context, limit int) ([]models.Job, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := js.pool.Query(ctx,
		`SELECT id, queue, payload, status, attempts, max_attempts, run_after, error_message, dead_at, created_at, updated_at
		 FROM jobs WHERE status IN ('queued','retrying') AND run_after <= NOW() ORDER BY run_after LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanJobs(rows)
}

// ListDead returns dead jobs.
func (js *JobStore) ListDead(ctx context.Context, limit int) ([]models.Job, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := js.pool.Query(ctx,
		`SELECT id, queue, payload, status, attempts, max_attempts, run_after, error_message, dead_at, created_at, updated_at
		 FROM jobs WHERE status='dead' ORDER BY dead_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanJobs(rows)
}

func scanJobs(rows pgx.Rows) ([]models.Job, error) {
	var jobs []models.Job
	for rows.Next() {
		var j models.Job
		var payload []byte
		err := rows.Scan(&j.ID, &j.Queue, &payload, &j.Status, &j.Attempts, &j.MaxAttempts, &j.RunAfter, &j.ErrorMessage, &j.DeadAt, &j.CreatedAt, &j.UpdatedAt)
		if err != nil {
			return nil, err
		}
		_ = json.Unmarshal(payload, &j.Payload)
		jobs = append(jobs, j)
	}
	return jobs, rows.Err()
}
