package mongo

import (
	"context"
	"time"

	"github.com/sukhera/uptime-monitor/internal/domain/service"
	"github.com/sukhera/uptime-monitor/internal/shared/errors"
	"github.com/sukhera/uptime-monitor/internal/shared/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// StatusLogRepository implements the status log repository interface for MongoDB
type StatusLogRepository struct {
	db Interface
}

// NewStatusLogRepository creates a new status log repository
func NewStatusLogRepository(db Interface) *StatusLogRepository {
	return &StatusLogRepository{
		db: db,
	}
}

// Create creates a new status log entry
func (r *StatusLogRepository) Create(ctx context.Context, log *service.StatusLog) error {
	if log == nil {
		return errors.NewValidationError("status log cannot be nil")
	}

	if log.ServiceName == "" {
		return errors.NewValidationError("service name is required")
	}

	if log.Status == "" {
		return errors.NewValidationError("status is required")
	}

	if r.db == nil {
		return errors.New("database connection is nil", errors.ErrorKindInternal)
	}

	// Set timestamp if not provided
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now().UTC()
	}

	_, err := r.db.StatusLogsCollection().InsertOne(ctx, log)
	if err != nil {
		return errors.NewWithCause("failed to create status log", errors.ErrorKindInternal, err)
	}

	return nil
}

// FindByServiceName retrieves status logs for a specific service
func (r *StatusLogRepository) FindByServiceName(ctx context.Context, serviceName string, limit int) ([]*service.StatusLog, error) {
	if serviceName == "" {
		return nil, errors.NewValidationError("service name is required")
	}

	if r.db == nil {
		return nil, errors.New("database connection is nil", errors.ErrorKindInternal)
	}

	filter := bson.M{"service_name": serviceName}
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}). // Most recent first
		SetLimit(int64(limit))

	cursor, err := r.db.StatusLogsCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, errors.NewWithCause("failed to find status logs for service", errors.ErrorKindInternal, err)
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log := logger.Get()
			log.Error(ctx, "Error closing cursor", err, nil)
		}
	}()

	var statusLogs []*service.StatusLog
	if err = cursor.All(ctx, &statusLogs); err != nil {
		return nil, errors.NewWithCause("failed to decode status logs", errors.ErrorKindInternal, err)
	}

	return statusLogs, nil
}

// FindRecent retrieves the most recent status logs across all services
func (r *StatusLogRepository) FindRecent(ctx context.Context, limit int) ([]*service.StatusLog, error) {
	if r.db == nil {
		return nil, errors.New("database connection is nil", errors.ErrorKindInternal)
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}). // Most recent first
		SetLimit(int64(limit))

	cursor, err := r.db.StatusLogsCollection().Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, errors.NewWithCause("failed to find recent status logs", errors.ErrorKindInternal, err)
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log := logger.Get()
			log.Error(ctx, "Error closing cursor", err, nil)
		}
	}()

	var statusLogs []*service.StatusLog
	if err = cursor.All(ctx, &statusLogs); err != nil {
		return nil, errors.NewWithCause("failed to decode status logs", errors.ErrorKindInternal, err)
	}

	return statusLogs, nil
}

// DeleteOlderThan removes status logs older than the specified cutoff time
func (r *StatusLogRepository) DeleteOlderThan(ctx context.Context, cutoff time.Time) error {
	if r.db == nil {
		return errors.New("database connection is nil", errors.ErrorKindInternal)
	}

	filter := bson.M{"timestamp": bson.M{"$lt": cutoff}}

	result, err := r.db.StatusLogsCollection().DeleteMany(ctx, filter)
	if err != nil {
		return errors.NewWithCause("failed to delete old status logs", errors.ErrorKindInternal, err)
	}

	if result.DeletedCount > 0 {
		log := logger.Get()
		log.Info(ctx, "Cleaned up old status logs", logger.Fields{
			"deleted_count": result.DeletedCount,
			"cutoff_date":   cutoff,
		})
	}

	return nil
}

// GetLatestByService retrieves the latest status log for each service
func (r *StatusLogRepository) GetLatestByService(ctx context.Context) (map[string]*service.StatusLog, error) {
	if r.db == nil {
		return nil, errors.New("database connection is nil", errors.ErrorKindInternal)
	}

	// MongoDB aggregation pipeline to get the latest status for each service
	pipeline := []bson.M{
		{
			"$sort": bson.M{"timestamp": -1},
		},
		{
			"$group": bson.M{
				"_id":        "$service_name",
				"latest_log": bson.M{"$first": "$$ROOT"},
			},
		},
	}

	cursor, err := r.db.StatusLogsCollection().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.NewWithCause("failed to aggregate latest status by service", errors.ErrorKindInternal, err)
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log := logger.Get()
			log.Error(ctx, "Error closing cursor", err, nil)
		}
	}()

	result := make(map[string]*service.StatusLog)

	for cursor.Next(ctx) {
		var doc struct {
			ID        string            `bson:"_id"`
			LatestLog service.StatusLog `bson:"latest_log"`
		}

		if err := cursor.Decode(&doc); err != nil {
			log := logger.Get()
			log.Error(ctx, "Error decoding aggregation result", err, nil)
			continue
		}

		result[doc.ID] = &doc.LatestLog
	}

	if err := cursor.Err(); err != nil {
		return nil, errors.NewWithCause("cursor error during aggregation", errors.ErrorKindInternal, err)
	}

	return result, nil
}
