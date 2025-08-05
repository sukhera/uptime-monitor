package mongo

import (
	"context"
	"fmt"

	"github.com/sukhera/uptime-monitor/internal/domain/service"
	"github.com/sukhera/uptime-monitor/internal/shared/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ServiceRepository implements the service repository interface for MongoDB
type ServiceRepository struct {
	db Interface
}

// NewServiceRepository creates a new service repository
func NewServiceRepository(db Interface) *ServiceRepository {
	return &ServiceRepository{
		db: db,
	}
}

// Create creates a new service
func (r *ServiceRepository) Create(ctx context.Context, svc *service.Service) error {
	if err := svc.Validate(); err != nil {
		return errors.NewValidationError(fmt.Sprintf("invalid service: %v", err))
	}

	result, err := r.db.ServicesCollection().InsertOne(ctx, svc)
	if err != nil {
		return errors.NewWithCause("failed to create service", errors.ErrorKindInternal, err)
	}

	// Note: Service struct doesn't have an ID field, so we don't set it
	_ = result

	return nil
}

// GetByID retrieves a service by its ID
func (r *ServiceRepository) GetByID(ctx context.Context, id string) (*service.Service, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid service ID format")
	}

	filter := bson.M{"_id": objectID}
	var svc service.Service
	err = r.db.ServicesCollection().FindOne(ctx, filter).Decode(&svc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("service not found")
		}
		return nil, errors.NewWithCause("failed to find service", errors.ErrorKindInternal, err)
	}

	return &svc, nil
}

// GetBySlug retrieves a service by slug
func (r *ServiceRepository) GetBySlug(ctx context.Context, slug string) (*service.Service, error) {
	filter := bson.M{"slug": slug}
	var svc service.Service
	err := r.db.ServicesCollection().FindOne(ctx, filter).Decode(&svc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("service not found")
		}
		return nil, errors.NewWithCause("failed to find service", errors.ErrorKindInternal, err)
	}

	return &svc, nil
}

// GetAll retrieves all services
func (r *ServiceRepository) GetAll(ctx context.Context) ([]*service.Service, error) {
	cursor, err := r.db.ServicesCollection().Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.NewWithCause("failed to find services", errors.ErrorKindInternal, err)
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Error closing cursor: %v\n", err)
		}
	}()

	var services []*service.Service
	if err = cursor.All(ctx, &services); err != nil {
		return nil, errors.NewWithCause("failed to decode services", errors.ErrorKindInternal, err)
	}

	return services, nil
}

// GetEnabled retrieves all enabled services
func (r *ServiceRepository) GetEnabled(ctx context.Context) ([]*service.Service, error) {
	filter := bson.M{"enabled": true}
	cursor, err := r.db.ServicesCollection().Find(ctx, filter)
	if err != nil {
		return nil, errors.NewWithCause("failed to find enabled services", errors.ErrorKindInternal, err)
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Error closing cursor: %v\n", err)
		}
	}()

	var services []*service.Service
	if err = cursor.All(ctx, &services); err != nil {
		return nil, errors.NewWithCause("failed to decode services", errors.ErrorKindInternal, err)
	}

	return services, nil
}

// Update updates an existing service
func (r *ServiceRepository) Update(ctx context.Context, svc *service.Service) error {
	if err := svc.Validate(); err != nil {
		return errors.NewValidationError(fmt.Sprintf("invalid service: %v", err))
	}

	// For now, we'll update by slug since Service doesn't have an ID field
	filter := bson.M{"slug": svc.Slug}
	update := bson.M{"$set": svc}

	result, err := r.db.ServicesCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.NewWithCause("failed to update service", errors.ErrorKindInternal, err)
	}

	if result.MatchedCount == 0 {
		return errors.NewNotFoundError("service not found")
	}

	return nil
}

// Delete deletes a service by its ID
func (r *ServiceRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.NewValidationError("invalid service ID format")
	}

	filter := bson.M{"_id": objectID}
	result, err := r.db.ServicesCollection().DeleteOne(ctx, filter)
	if err != nil {
		return errors.NewWithCause("failed to delete service", errors.ErrorKindInternal, err)
	}

	if result.DeletedCount == 0 {
		return errors.NewNotFoundError("service not found")
	}

	return nil
}

// SaveStatusLog saves a status log entry
func (r *ServiceRepository) SaveStatusLog(ctx context.Context, log *service.StatusLog) error {
	result, err := r.db.StatusLogsCollection().InsertOne(ctx, log)
	if err != nil {
		return errors.NewWithCause("failed to create status log", errors.ErrorKindInternal, err)
	}

	// Note: StatusLog struct doesn't have an ID field, so we don't set it
	_ = result

	return nil
}

// GetLatestStatus retrieves the latest status for all services
func (r *ServiceRepository) GetLatestStatus(ctx context.Context) ([]*service.ServiceStatus, error) {
	// This is a simplified implementation
	// In a real application, you might want to aggregate the latest status for each service
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}).SetLimit(100)
	cursor, err := r.db.StatusLogsCollection().Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, errors.NewWithCause("failed to find latest status logs", errors.ErrorKindInternal, err)
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Error closing cursor: %v\n", err)
		}
	}()

	var statusLogs []service.StatusLog
	if err = cursor.All(ctx, &statusLogs); err != nil {
		return nil, errors.NewWithCause("failed to decode status logs", errors.ErrorKindInternal, err)
	}

	// Convert to ServiceStatus (simplified)
	var statuses []*service.ServiceStatus
	for _, log := range statusLogs {
		status := &service.ServiceStatus{
			Name:      log.ServiceName,
			Status:    log.Status,
			Latency:   log.Latency,
			UpdatedAt: log.Timestamp,
			Error:     log.Error,
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}

// GetStatusHistory retrieves status history for a service
func (r *ServiceRepository) GetStatusHistory(ctx context.Context, serviceName string, limit int) ([]*service.StatusLog, error) {
	filter := bson.M{"service_name": serviceName}
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}).SetLimit(int64(limit))

	cursor, err := r.db.StatusLogsCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, errors.NewWithCause("failed to find status logs for service", errors.ErrorKindInternal, err)
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Error closing cursor: %v\n", err)
		}
	}()

	var statusLogs []*service.StatusLog
	if err = cursor.All(ctx, &statusLogs); err != nil {
		return nil, errors.NewWithCause("failed to decode status logs", errors.ErrorKindInternal, err)
	}

	return statusLogs, nil
}
