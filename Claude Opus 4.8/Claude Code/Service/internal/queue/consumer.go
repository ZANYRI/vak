package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/example/billing-service/internal/observability"
	"github.com/redis/go-redis/v9"
)

// Handler processes a job. Returning an error triggers retry/backoff/DLQ.
type Handler func(ctx context.Context, job Job) error

// Consumer reads jobs from the Redis stream and dispatches to handlers.
type Consumer struct {
	q           *Queue
	log         *slog.Logger
	metrics     *observability.Metrics
	consumerID  string
	handlers    map[string]Handler
	jobTimeout  time.Duration
}

func NewConsumer(q *Queue, log *slog.Logger, m *observability.Metrics, consumerID string) *Consumer {
	return &Consumer{
		q:          q,
		log:        log,
		metrics:    m,
		consumerID: consumerID,
		handlers:   make(map[string]Handler),
		jobTimeout: 30 * time.Second,
	}
}

// Register binds a handler to a job type.
func (c *Consumer) Register(jobType string, h Handler) { c.handlers[jobType] = h }

// Run consumes until ctx is cancelled (graceful shutdown).
func (c *Consumer) Run(ctx context.Context) error {
	c.log.Info("worker started", "consumer", c.consumerID)
	for {
		select {
		case <-ctx.Done():
			c.log.Info("worker shutting down", "consumer", c.consumerID)
			return nil
		default:
		}

		streams, err := c.q.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    groupName,
			Consumer: c.consumerID,
			Streams:  []string{streamMain, ">"},
			Count:    10,
			Block:    2 * time.Second,
		}).Result()
		if err != nil {
			if err == redis.Nil || ctx.Err() != nil {
				continue
			}
			c.log.Error("xreadgroup failed", "error", err.Error())
			time.Sleep(time.Second)
			continue
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				c.handleMessage(ctx, msg)
			}
		}
	}
}

func (c *Consumer) handleMessage(ctx context.Context, msg redis.XMessage) {
	raw, _ := msg.Values["job"].(string)
	var job Job
	if err := json.Unmarshal([]byte(raw), &job); err != nil {
		c.log.Error("invalid job payload, acking", "id", msg.ID, "error", err.Error())
		c.q.rdb.XAck(ctx, streamMain, groupName, msg.ID)
		return
	}

	// Idempotency: skip already-completed jobs.
	if done, _ := c.q.jobExistsCompleted(ctx, job.ID); done {
		c.q.rdb.XAck(ctx, streamMain, groupName, msg.ID)
		return
	}

	handler, ok := c.handlers[job.Type]
	if !ok {
		c.log.Warn("no handler for job type, sending to dead", "type", job.Type)
		c.toDead(ctx, msg.ID, job, "no handler registered")
		return
	}

	job.Attempts++
	c.q.markStatus(ctx, job.ID, "running", "", job.Attempts)

	err := c.runWithRecovery(ctx, handler, job)
	if err == nil {
		c.q.markStatus(ctx, job.ID, "completed", "", job.Attempts)
		c.q.rdb.XAck(ctx, streamMain, groupName, msg.ID)
		c.metrics.JobsProcessedTotal.WithLabelValues(job.Type).Inc()
		return
	}

	c.metrics.JobsFailedTotal.WithLabelValues(job.Type).Inc()
	c.log.Error("job failed", "type", job.Type, "attempt", job.Attempts, "error", err.Error())

	if job.Attempts >= job.MaxAttempts {
		c.toDead(ctx, msg.ID, job, err.Error())
		return
	}

	// Retry with exponential backoff: ack the current delivery and re-enqueue later.
	c.q.markStatus(ctx, job.ID, "retrying", err.Error(), job.Attempts)
	c.q.rdb.XAck(ctx, streamMain, groupName, msg.ID)
	delay := backoff(job.Attempts)
	go c.requeueAfter(ctx, job, delay)
}

func (c *Consumer) runWithRecovery(ctx context.Context, h Handler, job Job) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("panic in handler: %v", rec)
		}
	}()
	jobCtx, cancel := context.WithTimeout(ctx, c.jobTimeout)
	defer cancel()
	return h(jobCtx, job)
}

func (c *Consumer) requeueAfter(ctx context.Context, job Job, delay time.Duration) {
	select {
	case <-time.After(delay):
	case <-ctx.Done():
		return
	}
	if err := c.q.xadd(context.Background(), streamMain, job); err != nil {
		c.log.Error("failed to requeue job", "type", job.Type, "error", err.Error())
	}
}

func (c *Consumer) toDead(ctx context.Context, msgID string, job Job, reason string) {
	c.q.markStatus(ctx, job.ID, "dead", reason, job.Attempts)
	if err := c.q.xadd(ctx, streamDead, job); err != nil {
		c.log.Error("failed to write dead letter", "type", job.Type, "error", err.Error())
	}
	c.q.rdb.XAck(ctx, streamMain, groupName, msgID)
}
