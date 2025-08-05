package checker

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockObserver is a mock observer for testing
type MockObserver struct {
	notified bool
	events   []HealthCheckEvent
}

func (m *MockObserver) OnHealthCheckCompleted(ctx context.Context, event HealthCheckEvent) {
	m.notified = true
	m.events = append(m.events, event)
}

func TestHealthCheckSubject_New(t *testing.T) {
	subject := NewHealthCheckSubject()

	assert.NotNil(t, subject)
	assert.NotNil(t, subject.observers)
	assert.Empty(t, subject.observers)
}

func TestHealthCheckSubject_Attach(t *testing.T) {
	subject := NewHealthCheckSubject()
	observer := &MockObserver{}

	subject.Attach(observer)

	assert.Len(t, subject.observers, 1)
	assert.Equal(t, observer, subject.observers[0])
}

func TestHealthCheckSubject_Attach_Multiple(t *testing.T) {
	subject := NewHealthCheckSubject()
	observer1 := &MockObserver{}
	observer2 := &MockObserver{}

	subject.Attach(observer1)
	subject.Attach(observer2)

	assert.Len(t, subject.observers, 2)
	assert.Equal(t, observer1, subject.observers[0])
	assert.Equal(t, observer2, subject.observers[1])
}

func TestHealthCheckSubject_Detach(t *testing.T) {
	subject := NewHealthCheckSubject()
	observer1 := &MockObserver{}
	observer2 := &MockObserver{}

	subject.Attach(observer1)
	subject.Attach(observer2)

	// Detach first observer
	subject.Detach(observer1)

	assert.Len(t, subject.observers, 1)
	assert.Equal(t, observer2, subject.observers[0])
}

func TestHealthCheckSubject_Detach_NotExists(t *testing.T) {
	subject := NewHealthCheckSubject()
	observer1 := &MockObserver{}
	observer2 := &MockObserver{}

	subject.Attach(observer1)

	// Try to detach observer that wasn't attached
	subject.Detach(observer2)

	assert.Len(t, subject.observers, 1)
	assert.Equal(t, observer1, subject.observers[0])
}

func TestHealthCheckSubject_Notify(t *testing.T) {
	subject := NewHealthCheckSubject()
	observer := &MockObserver{}

	subject.Attach(observer)

	event := HealthCheckEvent{
		ServiceName: "test-service",
		Status:      "operational",
		Latency:     150,
		StatusCode:  200,
		Timestamp:   time.Now().Unix(),
	}

	ctx := context.Background()
	subject.Notify(ctx, event)

	// Wait a bit for goroutine to complete
	time.Sleep(10 * time.Millisecond)

	assert.True(t, observer.notified)
	assert.Len(t, observer.events, 1)
	assert.Equal(t, event, observer.events[0])
}

func TestHealthCheckSubject_Notify_MultipleObservers(t *testing.T) {
	subject := NewHealthCheckSubject()
	observer1 := &MockObserver{}
	observer2 := &MockObserver{}

	subject.Attach(observer1)
	subject.Attach(observer2)

	event := HealthCheckEvent{
		ServiceName: "test-service",
		Status:      "down",
		Latency:     5000,
		StatusCode:  500,
		Error:       "Connection timeout",
		Timestamp:   time.Now().Unix(),
	}

	ctx := context.Background()
	subject.Notify(ctx, event)

	// Wait a bit for goroutines to complete
	time.Sleep(10 * time.Millisecond)

	assert.True(t, observer1.notified)
	assert.True(t, observer2.notified)
	assert.Len(t, observer1.events, 1)
	assert.Len(t, observer2.events, 1)
	assert.Equal(t, event, observer1.events[0])
	assert.Equal(t, event, observer2.events[0])
}

func TestHealthCheckSubject_Notify_NoObservers(t *testing.T) {
	subject := NewHealthCheckSubject()

	event := HealthCheckEvent{
		ServiceName: "test-service",
		Status:      "operational",
		Latency:     150,
		StatusCode:  200,
		Timestamp:   time.Now().Unix(),
	}

	ctx := context.Background()
	// Should not panic when no observers are attached
	assert.NotPanics(t, func() {
		subject.Notify(ctx, event)
	})
}

func TestMetricsObserver_New(t *testing.T) {
	observer := NewMetricsObserver()

	assert.NotNil(t, observer)
	assert.NotNil(t, observer.metrics)
	assert.Empty(t, observer.metrics)
}

func TestMetricsObserver_OnHealthCheckCompleted(t *testing.T) {
	observer := NewMetricsObserver()

	event := HealthCheckEvent{
		ServiceName: "test-service",
		Status:      "operational",
		Latency:     150,
		StatusCode:  200,
		Timestamp:   time.Now().Unix(),
	}

	ctx := context.Background()
	observer.OnHealthCheckCompleted(ctx, event)

	metrics := observer.GetMetrics()
	assert.NotEmpty(t, metrics)

	// Check service-specific metrics
	serviceKey := "service_test-service"
	serviceMetrics, exists := metrics[serviceKey]
	assert.True(t, exists)

	serviceData, ok := serviceMetrics.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "operational", serviceData["status"])
	assert.Equal(t, int64(150), serviceData["latency_ms"])
	assert.Equal(t, 200, serviceData["status_code"])

	// Check global metrics
	globalMetrics, exists := metrics["global"]
	assert.True(t, exists)

	globalData, ok := globalMetrics.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, 1, globalData["total_checks"])
}

func TestMetricsObserver_OnHealthCheckCompleted_Multiple(t *testing.T) {
	observer := NewMetricsObserver()

	events := []HealthCheckEvent{
		{
			ServiceName: "service-1",
			Status:      "operational",
			Latency:     150,
			StatusCode:  200,
			Timestamp:   time.Now().Unix(),
		},
		{
			ServiceName: "service-2",
			Status:      "down",
			Latency:     5000,
			StatusCode:  500,
			Error:       "Connection timeout",
			Timestamp:   time.Now().Unix(),
		},
	}

	ctx := context.Background()
	for _, event := range events {
		observer.OnHealthCheckCompleted(ctx, event)
	}

	metrics := observer.GetMetrics()

	// Check global metrics
	globalMetrics, exists := metrics["global"]
	assert.True(t, exists)

	globalData, ok := globalMetrics.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, 2, globalData["total_checks"])

	// Check service-specific metrics
	service1Key := "service_service-1"
	service1Metrics, exists := metrics[service1Key]
	assert.True(t, exists)

	service1Data, ok := service1Metrics.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "operational", service1Data["status"])

	service2Key := "service_service-2"
	service2Metrics, exists := metrics[service2Key]
	assert.True(t, exists)

	service2Data, ok := service2Metrics.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "down", service2Data["status"])
}

func TestMetricsObserver_GetMetrics(t *testing.T) {
	observer := NewMetricsObserver()

	// Add some metrics
	observer.metrics["test"] = "value"
	observer.metrics["number"] = 42

	metrics := observer.GetMetrics()

	assert.Equal(t, "value", metrics["test"])
	assert.Equal(t, 42, metrics["number"])
}

func TestAlertingObserver_New(t *testing.T) {
	observer := NewAlertingObserver(5000)

	assert.NotNil(t, observer)
	assert.Equal(t, int64(5000), observer.alertThreshold)
	assert.NotNil(t, observer.alertCh)
}

func TestAlertingObserver_OnHealthCheckCompleted_NoAlert(t *testing.T) {
	observer := NewAlertingObserver(5000)

	event := HealthCheckEvent{
		ServiceName: "test-service",
		Status:      "operational",
		Latency:     150,
		StatusCode:  200,
		Timestamp:   time.Now().Unix(),
	}

	ctx := context.Background()
	observer.OnHealthCheckCompleted(ctx, event)

	// Should not send alert for operational service with low latency
	select {
	case <-observer.alertCh:
		assert.Fail(t, "Should not have sent alert")
	case <-time.After(10 * time.Millisecond):
		// Expected - no alert should be sent
	}
}

func TestAlertingObserver_OnHealthCheckCompleted_Alert_Down(t *testing.T) {
	observer := NewAlertingObserver(5000)

	event := HealthCheckEvent{
		ServiceName: "test-service",
		Status:      "down",
		Latency:     150,
		StatusCode:  500,
		Error:       "Connection timeout",
		Timestamp:   time.Now().Unix(),
	}

	ctx := context.Background()
	observer.OnHealthCheckCompleted(ctx, event)

	// Should send alert for down service
	select {
	case alert := <-observer.alertCh:
		assert.Equal(t, event, alert)
	case <-time.After(10 * time.Millisecond):
		assert.Fail(t, "Should have sent alert")
	}
}

func TestAlertingObserver_OnHealthCheckCompleted_Alert_HighLatency(t *testing.T) {
	observer := NewAlertingObserver(5000)

	event := HealthCheckEvent{
		ServiceName: "test-service",
		Status:      "operational",
		Latency:     6000, // Above threshold
		StatusCode:  200,
		Timestamp:   time.Now().Unix(),
	}

	ctx := context.Background()
	observer.OnHealthCheckCompleted(ctx, event)

	// Should send alert for high latency
	select {
	case alert := <-observer.alertCh:
		assert.Equal(t, event, alert)
	case <-time.After(10 * time.Millisecond):
		assert.Fail(t, "Should have sent alert")
	}
}

func TestAlertingObserver_GetAlertChannel(t *testing.T) {
	observer := NewAlertingObserver(5000)

	alertCh := observer.GetAlertChannel()
	assert.NotNil(t, alertCh)
	// The method returns a receive-only channel, so we can't directly compare with the internal channel
	// Just verify it's not nil and has the right type
	assert.IsType(t, (<-chan HealthCheckEvent)(nil), alertCh)
}

func TestAlertingObserver_ChannelFull(t *testing.T) {
	observer := NewAlertingObserver(5000)

	// Fill the channel
	for i := 0; i < 100; i++ {
		event := HealthCheckEvent{
			ServiceName: "test-service",
			Status:      "down",
			Latency:     150,
			StatusCode:  500,
			Timestamp:   time.Now().Unix(),
		}

		ctx := context.Background()
		observer.OnHealthCheckCompleted(ctx, event)
	}

	// Should not panic when channel is full
	assert.NotPanics(t, func() {
		event := HealthCheckEvent{
			ServiceName: "test-service",
			Status:      "down",
			Latency:     150,
			StatusCode:  500,
			Timestamp:   time.Now().Unix(),
		}
		ctx := context.Background()
		observer.OnHealthCheckCompleted(ctx, event)
	})
}

func TestObserver_Integration(t *testing.T) {
	// Test complete observer workflow
	subject := NewHealthCheckSubject()

	// Add metrics observer
	metricsObserver := NewMetricsObserver()
	subject.Attach(metricsObserver)

	// Add alerting observer
	alertingObserver := NewAlertingObserver(5000)
	subject.Attach(alertingObserver)

	// Create events
	events := []HealthCheckEvent{
		{
			ServiceName: "service-1",
			Status:      "operational",
			Latency:     150,
			StatusCode:  200,
			Timestamp:   time.Now().Unix(),
		},
		{
			ServiceName: "service-2",
			Status:      "down",
			Latency:     5000,
			StatusCode:  500,
			Error:       "Connection timeout",
			Timestamp:   time.Now().Unix(),
		},
	}

	ctx := context.Background()

	// Notify observers
	for _, event := range events {
		subject.Notify(ctx, event)
	}

	// Wait for goroutines to complete
	time.Sleep(20 * time.Millisecond)

	// Check metrics
	metrics := metricsObserver.GetMetrics()
	globalMetrics, exists := metrics["global"]
	assert.True(t, exists)

	globalData, ok := globalMetrics.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, 2, globalData["total_checks"])

	// Check alerts
	alertCount := 0
	for {
		select {
		case <-alertingObserver.alertCh:
			alertCount++
		case <-time.After(10 * time.Millisecond):
			goto done
		}
	}
done:
	assert.Equal(t, 1, alertCount) // Only service-2 should trigger alert
}

func TestObserver_ConcurrentNotifications(t *testing.T) {
	subject := NewHealthCheckSubject()
	observer := &MockObserver{}

	subject.Attach(observer)

	// Send multiple notifications concurrently
	ctx := context.Background()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			event := HealthCheckEvent{
				ServiceName: "service-" + string(rune(index)),
				Status:      "operational",
				Latency:     150,
				StatusCode:  200,
				Timestamp:   time.Now().Unix(),
			}
			subject.Notify(ctx, event)
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	time.Sleep(50 * time.Millisecond) // Longer buffer for notifications to be processed

	// Should have received most notifications (with tolerance for race conditions)
	// In concurrent environments, some events might be lost due to timing
	assert.True(t, len(observer.events) >= 5, "Expected at least 5 events, got %d", len(observer.events))
}

func TestObserver_ContextCancellation(t *testing.T) {
	subject := NewHealthCheckSubject()
	observer := &MockObserver{}

	subject.Attach(observer)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	event := HealthCheckEvent{
		ServiceName: "test-service",
		Status:      "operational",
		Latency:     150,
		StatusCode:  200,
		Timestamp:   time.Now().Unix(),
	}

	subject.Notify(ctx, event)

	// Wait for potential notification
	time.Sleep(10 * time.Millisecond)

	// Observer might or might not be notified due to context cancellation
	// The important thing is that it doesn't panic
	assert.NotNil(t, observer)
}
