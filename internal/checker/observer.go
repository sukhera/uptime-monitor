package checker

import (
	"context"
	"fmt"
	"sync"

	"github.com/sukhera/uptime-monitor/internal/shared/logger"
)

// HealthCheckEvent represents a health check event
type HealthCheckEvent struct {
	ServiceName string
	Status      string
	Latency     int64
	StatusCode  int
	Error       string
	Timestamp   int64
}

// HealthCheckObserver defines the interface for health check observers
type HealthCheckObserver interface {
	OnHealthCheckCompleted(ctx context.Context, event HealthCheckEvent)
}

// HealthCheckSubject manages observers and notifies them of events
type HealthCheckSubject struct {
	observers []HealthCheckObserver
	mu        sync.RWMutex
}

// NewHealthCheckSubject creates a new health check subject
func NewHealthCheckSubject() *HealthCheckSubject {
	return &HealthCheckSubject{
		observers: make([]HealthCheckObserver, 0),
	}
}

// Attach adds an observer to the subject
func (s *HealthCheckSubject) Attach(observer HealthCheckObserver) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.observers = append(s.observers, observer)
}

// Detach removes an observer from the subject
func (s *HealthCheckSubject) Detach(observer HealthCheckObserver) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, obs := range s.observers {
		if obs == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

// Notify notifies all observers of a health check event
func (s *HealthCheckSubject) Notify(ctx context.Context, event HealthCheckEvent) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, observer := range s.observers {
		go observer.OnHealthCheckCompleted(ctx, event)
	}
}

// LoggingObserver logs health check events
type LoggingObserver struct {
	logger logger.Logger
}

// NewLoggingObserver creates a new logging observer
func NewLoggingObserver(logger logger.Logger) *LoggingObserver {
	return &LoggingObserver{
		logger: logger,
	}
}

// OnHealthCheckCompleted logs health check events
func (o *LoggingObserver) OnHealthCheckCompleted(ctx context.Context, event HealthCheckEvent) {
	fields := logger.Fields{
		"service_name": event.ServiceName,
		"status":       event.Status,
		"latency_ms":   event.Latency,
		"status_code":  event.StatusCode,
	}

	if event.Error != "" {
		fields["error"] = event.Error
		o.logger.Error(ctx, "Health check failed", fmt.Errorf("%s", event.Error), fields)
	} else {
		o.logger.Info(ctx, "Health check completed", fields)
	}
}

// MetricsObserver collects metrics from health check events
type MetricsObserver struct {
	metrics map[string]interface{}
	mu      sync.RWMutex
}

// NewMetricsObserver creates a new metrics observer
func NewMetricsObserver() *MetricsObserver {
	return &MetricsObserver{
		metrics: make(map[string]interface{}),
	}
}

// OnHealthCheckCompleted collects metrics from health check events
func (o *MetricsObserver) OnHealthCheckCompleted(ctx context.Context, event HealthCheckEvent) {
	o.mu.Lock()
	defer o.mu.Unlock()

	// Update service-specific metrics
	serviceKey := "service_" + event.ServiceName
	o.metrics[serviceKey] = map[string]interface{}{
		"status":      event.Status,
		"latency_ms":  event.Latency,
		"status_code": event.StatusCode,
		"timestamp":   event.Timestamp,
	}

	// Update global metrics
	if globalMetrics, exists := o.metrics["global"]; exists {
		if global, ok := globalMetrics.(map[string]interface{}); ok {
			if totalChecks, exists := global["total_checks"]; exists {
				if total, ok := totalChecks.(int); ok {
					global["total_checks"] = total + 1
				}
			} else {
				global["total_checks"] = 1
			}
		}
	} else {
		o.metrics["global"] = map[string]interface{}{
			"total_checks": 1,
		}
	}
}

// GetMetrics returns the collected metrics
func (o *MetricsObserver) GetMetrics() map[string]interface{} {
	o.mu.RLock()
	defer o.mu.RUnlock()

	metrics := make(map[string]interface{})
	for k, v := range o.metrics {
		metrics[k] = v
	}
	return metrics
}

// AlertingObserver handles alerting based on health check events
type AlertingObserver struct {
	alertThreshold int64 // latency threshold in milliseconds
	alertCh        chan HealthCheckEvent
}

// NewAlertingObserver creates a new alerting observer
func NewAlertingObserver(alertThreshold int64) *AlertingObserver {
	return &AlertingObserver{
		alertThreshold: alertThreshold,
		alertCh:        make(chan HealthCheckEvent, 100),
	}
}

// OnHealthCheckCompleted checks if an alert should be triggered
func (o *AlertingObserver) OnHealthCheckCompleted(ctx context.Context, event HealthCheckEvent) {
	// Check if service is down or latency is too high
	if event.Status == "down" || event.Latency > o.alertThreshold {
		select {
		case o.alertCh <- event:
			// Alert sent successfully
		default:
			// Alert channel is full, log the dropped alert
		}
	}
}

// GetAlertChannel returns the alert channel
func (o *AlertingObserver) GetAlertChannel() <-chan HealthCheckEvent {
	return o.alertCh
}
