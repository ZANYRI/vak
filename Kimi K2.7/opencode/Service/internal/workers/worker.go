package workers

import (
	"fmt"
	"runtime/debug"

	"billing-service/internal/observability"
	"billing-service/internal/queue"
	"go.uber.org/zap"
)

// Worker consumes jobs from the queue and executes registered handlers.
type Worker struct {
	registry *Registry
	client   *queue.Client
	logger   *zap.Logger
}

// NewWorker creates a worker.
func NewWorker(registry *Registry, client *queue.Client, logger *zap.Logger) *Worker {
	return &Worker{registry: registry, client: client, logger: logger}
}

// Run starts consuming from all registered queues.
func (w *Worker) Run() error {
	for _, q := range w.registry.Queues() {
		queueName := q
		handler := w.registry.Handler(queueName)
		wrapped := func(payload map[string]interface{}) (err error) {
			defer func() {
				if r := recover(); r != nil {
					w.logger.Error("job panic recovered", zap.Any("panic", r), zap.String("stack", string(debug.Stack())))
					err = fmt.Errorf("panic: %v", r)
				}
			}()
			return handler(payload)
		}
		if err := w.client.Subscribe(queueName, wrapped); err != nil {
			return fmt.Errorf("subscribe to %s: %w", queueName, err)
		}
		w.logger.Info("subscribed to queue", zap.String("queue", queueName))
	}
	return nil
}

// Wait blocks indefinitely while the queues are being consumed.
func (w *Worker) Wait() {
	select {}
}

// Metrics returns observable metrics names.
func init() {
	observability.IncJobProcessed("all", "started")
}
