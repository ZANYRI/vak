package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

// Client wraps a NATS connection and JetStream context.
type Client struct {
	nc      *nats.Conn
	js      nats.JetStreamContext
	stream  string
	store   *JobStore
}

// Queue names.
const (
	QueueInvoiceGenerate   = "invoice.generate"
	QueueInvoiceFinalize   = "invoice.finalize"
	QueuePaymentProcess    = "payment.process"
	QueueSubscriptionRenew = "subscription.renew"
	QueueExpireTrial       = "subscription.expire_trial"
	QueueEmailInvoice      = "email.invoice_created"
	QueueEmailPaymentFail  = "email.payment_failed"
	QueueUsageAggregate    = "usage.aggregate"
)

// NewClient connects to NATS and ensures the stream exists.
func NewClient(url, stream string, store *JobStore) (*Client, error) {
	nc, err := nats.Connect(url, nats.Timeout(10*time.Second), nats.ReconnectWait(2*time.Second))
	if err != nil {
		return nil, fmt.Errorf("connect to nats: %w", err)
	}
	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("jetstream: %w", err)
	}
	c := &Client{nc: nc, js: js, stream: stream, store: store}
	if err := c.ensureStream(); err != nil {
		nc.Close()
		return nil, err
	}
	return c, nil
}

// Close closes the NATS connection.
func (c *Client) Close() {
	if c.nc != nil {
		c.nc.Close()
	}
}

// Publish enqueues a job both in the database and NATS JetStream.
func (c *Client) Publish(ctx context.Context, queue string, payload map[string]interface{}, maxAttempts int) (string, error) {
	id, err := c.store.CreateJob(ctx, queue, payload, maxAttempts)
	if err != nil {
		return "", err
	}
	msgPayload := map[string]interface{}{"job_id": id.String(), "queue": queue, "payload": payload}
	data, _ := json.Marshal(msgPayload)
	_, err = c.js.Publish(queue, data, nats.Context(ctx))
	if err != nil {
		// mark dead if cannot publish
		_ = c.store.MarkDead(ctx, id, fmt.Sprintf("publish error: %v", err))
		return "", fmt.Errorf("publish job: %w", err)
	}
	return id.String(), nil
}

// Subscribe consumes messages for a queue. The handler is passed the job payload.
func (c *Client) Subscribe(queue string, handler func(map[string]interface{}) error) error {
	durable := fmt.Sprintf("%s-consumer", strings.ReplaceAll(queue, ".", "-"))
	_, err := c.js.Subscribe(queue, func(msg *nats.Msg) {
		var envelope struct {
			JobID   string                 `json:"job_id"`
			Queue   string                 `json:"queue"`
			Payload map[string]interface{} `json:"payload"`
		}
		if err := json.Unmarshal(msg.Data, &envelope); err != nil {
			_ = msg.Ack()
			return
		}
		jobID, err := parseUUID(envelope.JobID)
		if err != nil {
			_ = msg.Ack()
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		if err := c.store.MarkRunning(ctx, jobID); err != nil {
			_ = msg.NakWithDelay(5 * time.Second)
			return
		}
		if err := handler(envelope.Payload); err != nil {
			_ = c.store.MarkFailed(ctx, jobID, err.Error())
			_ = msg.NakWithDelay(5 * time.Second)
			return
		}
		_ = c.store.MarkCompleted(ctx, jobID)
		_ = msg.Ack()
	}, nats.Durable(durable), nats.ManualAck(), nats.MaxDeliver(1))
	return err
}

// RequeueRetries re-publishes retryable jobs.
func (c *Client) RequeueRetries(ctx context.Context) error {
	jobs, err := c.store.ListRetryable(ctx, 100)
	if err != nil {
		return err
	}
	for _, job := range jobs {
		payload := map[string]interface{}{"job_id": job.ID.String(), "queue": job.Queue, "payload": job.Payload}
		data, _ := json.Marshal(payload)
		_, _ = c.js.Publish(job.Queue, data)
	}
	return nil
}

func (c *Client) ensureStream() error {
	_, err := c.js.StreamInfo(c.stream)
	if err == nil {
		return nil
	}
	_, err = c.js.AddStream(&nats.StreamConfig{
		Name:     c.stream,
		Subjects: []string{"invoice.*", "payment.*", "subscription.*", "email.*", "usage.*"},
		Storage:  nats.FileStorage,
		Replicas: 1,
	})
	return err
}

func parseUUID(s string) (UUID, error) {
	return parseUUIDInternal(s)
}
