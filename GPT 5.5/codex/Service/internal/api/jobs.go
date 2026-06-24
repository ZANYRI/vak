package api

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/example/billing-service/internal/queue"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (a *App) enqueue(ctx context.Context, typ string, payload any) error {
	raw, e := json.Marshal(payload)
	if e != nil {
		return e
	}
	jid := uuid.New()
	_, e = a.DB.Exec(ctx, `INSERT INTO jobs(id,type,payload,status) VALUES($1,$2,$3,'queued')`, jid, typ, raw)
	if e != nil {
		return e
	}
	if a.Queue == nil {
		return fmt.Errorf("queue unavailable")
	}
	if e = a.Queue.Publish(ctx, queue.Message{JobID: jid.String(), Type: typ, Payload: raw}); e != nil {
		_, _ = a.DB.Exec(ctx, `UPDATE jobs SET status='retrying',last_error=$2,updated_at=now() WHERE id=$1`, jid, e.Error())
		return e
	}
	return nil
}

// ProcessJob is deliberately small: API mutations are transactional; asynchronous work updates job state,
// and scheduled/notification jobs can be retried safely because each job has one persisted ID.
func (a *App) ProcessJob(ctx context.Context, m queue.Message) error {
	var attempts, maxAttempts int
	e := a.DB.QueryRow(ctx, `UPDATE jobs SET status='running',attempts=attempts+1,updated_at=now() WHERE id=$1 AND status IN ('queued','retrying') RETURNING attempts,max_attempts`, m.JobID).Scan(&attempts, &maxAttempts)
	if e == pgx.ErrNoRows {
		return nil
	} // an already completed/dead job is safe to acknowledge
	if e != nil {
		return e
	}
	switch m.Type {
	case "usage.aggregate":
		var p struct {
			SubscriptionID string `json:"subscription_id"`
		}
		if json.Unmarshal(m.Payload, &p) == nil && p.SubscriptionID != "" {
			_, e = a.DB.Exec(ctx, `INSERT INTO usage_summaries(id,subscription_id,metric,period_start,period_end,quantity) SELECT $1,subscription_id,metric,date_trunc('month',recorded_at),date_trunc('month',recorded_at)+interval '1 month',sum(quantity) FROM usage_events WHERE subscription_id=$2 GROUP BY subscription_id,metric,date_trunc('month',recorded_at) ON CONFLICT(subscription_id,metric,period_start,period_end) DO UPDATE SET quantity=EXCLUDED.quantity,updated_at=now()`, uuid.New(), p.SubscriptionID)
		}
	case "subscription.expire_trial":
		_, e = a.DB.Exec(ctx, `UPDATE subscriptions SET status='active',updated_at=now() WHERE status='trialing' AND trial_end <= now()`)
	case "subscription.renew":
		_, e = a.DB.Exec(ctx, `UPDATE subscriptions SET current_period_start=current_period_end,current_period_end=CASE WHEN p.billing_interval='yearly' THEN current_period_end+interval '1 year' ELSE current_period_end+interval '1 month' END,updated_at=now() FROM plans p WHERE subscriptions.plan_id=p.id AND subscriptions.status='active' AND subscriptions.current_period_end<=now()`)
	}
	if e != nil {
		a.metric.jobsFailed.Add(1)
		if attempts >= maxAttempts {
			_, _ = a.DB.Exec(ctx, `UPDATE jobs SET status='dead',last_error=$2,updated_at=now() WHERE id=$1`, m.JobID, e.Error())
			return nil // persistent job state is the dead-letter queue; acknowledge NATS delivery
		}
		_, _ = a.DB.Exec(ctx, `UPDATE jobs SET status='retrying',last_error=$2,updated_at=now() WHERE id=$1`, m.JobID, e.Error())
		return e
	}
	_, e = a.DB.Exec(ctx, `UPDATE jobs SET status='completed',updated_at=now() WHERE id=$1`, m.JobID)
	if e == nil {
		a.metric.jobsProcessed.Add(1)
	}
	return e
}
func (a *App) PublishScheduled(ctx context.Context) error {
	types := []string{"subscription.expire_trial", "subscription.renew", "usage.aggregate"}
	for _, t := range types {
		jobCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		e := a.enqueue(jobCtx, t, map[string]string{"scheduled_at": time.Now().UTC().Format(time.RFC3339)})
		cancel()
		if e != nil {
			return e
		}
	}
	return nil
}
