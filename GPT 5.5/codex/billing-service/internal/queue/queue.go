package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
)

const Stream = "BILLING"
const Subject = "billing.jobs"

type Message struct {
	JobID   string          `json:"job_id"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
	Attempt int             `json:"attempt"`
}
type Client struct {
	nc     *nats.Conn
	js     nats.JetStreamContext
	logger *slog.Logger
}

func Connect(url string, logger *slog.Logger) (*Client, error) {
	nc, e := nats.Connect(url, nats.Name("billing-service"), nats.Timeout(5*time.Second), nats.MaxReconnects(-1), nats.ReconnectWait(time.Second))
	if e != nil {
		return nil, e
	}
	js, e := nc.JetStream()
	if e != nil {
		nc.Close()
		return nil, e
	}
	if _, e = js.AddStream(&nats.StreamConfig{Name: Stream, Subjects: []string{"billing.>"}, Storage: nats.FileStorage, Retention: nats.WorkQueuePolicy, Discard: nats.DiscardOld, MaxAge: 7 * 24 * time.Hour}); e != nil && e != nats.ErrStreamNameAlreadyInUse {
		nc.Close()
		return nil, e
	}
	return &Client{nc, js, logger}, nil
}
func (c *Client) Close() {
	if c != nil && c.nc != nil {
		c.nc.Drain()
		c.nc.Close()
	}
}
func (c *Client) Publish(ctx context.Context, m Message) error {
	b, e := json.Marshal(m)
	if e != nil {
		return e
	}
	_, e = c.js.Publish(Subject, b, nats.Context(ctx))
	return e
}
func (c *Client) Consume(ctx context.Context, handler func(context.Context, Message) error) error {
	sub, e := c.js.PullSubscribe(Subject, "billing-worker", nats.BindStream(Stream))
	if e != nil && e != nats.ErrConsumerNameAlreadyInUse {
		return e
	}
	for {
		fetchCtx, cancel := context.WithTimeout(ctx, time.Second)
		msgs, e := sub.Fetch(1, nats.Context(fetchCtx))
		cancel()
		if e != nil {
			if e == nats.ErrTimeout || e == context.DeadlineExceeded {
				continue
			}
			if ctx.Err() != nil {
				return nil
			}
			return e
		}
		for _, msg := range msgs {
			var data Message
			if e := json.Unmarshal(msg.Data, &data); e != nil {
				_ = msg.Term()
				continue
			}
			if metadata, metaErr := msg.Metadata(); metaErr == nil && metadata.NumDelivered > 0 {
				data.Attempt = int(metadata.NumDelivered - 1)
			}
			jobCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			e = runSafely(jobCtx, handler, data)
			cancel()
			if e != nil {
				c.logger.Error("job failed", "job_id", data.JobID, "type", data.Type, "error", e)
				_ = msg.NakWithDelay(backoff(data.Attempt))
				continue
			}
			_ = msg.Ack()
		}
	}
}
func runSafely(ctx context.Context, h func(context.Context, Message) error, m Message) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	return h(ctx, m)
}
func backoff(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	if attempt > 8 {
		attempt = 8
	}
	return time.Second * time.Duration(1<<attempt)
}
