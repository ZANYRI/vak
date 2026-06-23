package audit

import (
	"context"

	"billing-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service logs audit events.
type Service struct {
	pool *pgxpool.Pool
}

// NewService creates an audit service.
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// LogEvent writes an audit log row.
func (s *Service) LogEvent(ctx context.Context, actorID *uuid.UUID, action, resourceType string, resourceID *uuid.UUID, metadata map[string]interface{}, ip, userAgent string) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO audit_logs (id, actor_user_id, action, resource_type, resource_id, metadata, ip_address, user_agent)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		uuid.New(), actorID, action, resourceType, resourceID, metadata, ip, userAgent)
	return err
}

// List returns audit logs with pagination.
func (s *Service) List(ctx context.Context, limit, offset int) ([]models.AuditLog, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	rows, err := s.pool.Query(ctx,
		`SELECT id, actor_user_id, action, resource_type, resource_id, metadata, ip_address, user_agent, created_at
		 FROM audit_logs ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []models.AuditLog
	for rows.Next() {
		var a models.AuditLog
		err := rows.Scan(&a.ID, &a.ActorUserID, &a.Action, &a.ResourceType, &a.ResourceID, &a.Metadata, &a.IPAddress, &a.UserAgent, &a.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, a)
	}
	var total int64
	_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM audit_logs`).Scan(&total)
	return list, total, rows.Err()
}
