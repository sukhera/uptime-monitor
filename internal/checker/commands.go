package checker

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sukhera/uptime-monitor/internal/domain/service"
	"github.com/sukhera/uptime-monitor/internal/shared/logger"
)

// HealthCheckCommand represents a health check command
type HealthCheckCommand interface {
	Execute(ctx context.Context) service.StatusLog
	GetServiceName() string
}

// HTTPHealthCheckCommand implements HTTP health checks
type HTTPHealthCheckCommand struct {
	service service.Service
	client  HTTPClient
}

// NewHTTPHealthCheckCommand creates a new HTTP health check command
func NewHTTPHealthCheckCommand(service service.Service, client HTTPClient) *HTTPHealthCheckCommand {
	return &HTTPHealthCheckCommand{
		service: service,
		client:  client,
	}
}

// Execute performs the HTTP health check
func (cmd *HTTPHealthCheckCommand) Execute(ctx context.Context) service.StatusLog {
	const maxRetries = 3
	const retryDelay = 500 * time.Millisecond

	req, err := http.NewRequestWithContext(ctx, "GET", cmd.service.URL, nil)
	if err != nil {
		return service.StatusLog{
			ServiceName: cmd.service.Name,
			Status:      statusDown,
			Latency:     0,
			StatusCode:  0,
			Error:       fmt.Errorf("Failed to create request: %w", err).Error(),
			Timestamp:   time.Now(),
		}
	}

	for k, v := range cmd.service.Headers {
		req.Header.Set(k, v)
	}

	var resp *http.Response
	var latency int64
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		start := time.Now()
		resp, err = cmd.client.Do(req)
		latency = time.Since(start).Milliseconds()

		if err == nil {
			break
		}

		lastErr = err
		if attempt < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	statusLog := service.StatusLog{
		ServiceName: cmd.service.Name,
		Latency:     latency,
		Timestamp:   time.Now(),
	}

	if lastErr != nil && resp == nil {
		statusLog.Status = statusDown
		statusLog.Error = fmt.Errorf("Request failed after %d attempts: %w", maxRetries, lastErr).Error()
		return statusLog
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Log error but don't fail the health check
			log := logger.Get()
			log.Error(ctx, "Error closing response body", err, nil)
		}
	}()

	statusLog.StatusCode = resp.StatusCode

	if resp.StatusCode == cmd.service.ExpectedStatus {
		statusLog.Status = statusOperational
	} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		statusLog.Status = statusDegraded
	} else {
		statusLog.Status = statusDown
		statusLog.Error = fmt.Sprintf("Unexpected status code: %d", resp.StatusCode)
	}

	return statusLog
}

// GetServiceName returns the service name
func (cmd *HTTPHealthCheckCommand) GetServiceName() string {
	return cmd.service.Name
}

// HealthCheckInvoker manages health check commands
type HealthCheckInvoker struct {
	commands []HealthCheckCommand
}

// NewHealthCheckInvoker creates a new health check invoker
func NewHealthCheckInvoker() *HealthCheckInvoker {
	return &HealthCheckInvoker{
		commands: make([]HealthCheckCommand, 0),
	}
}

// AddCommand adds a health check command
func (invoker *HealthCheckInvoker) AddCommand(command HealthCheckCommand) {
	invoker.commands = append(invoker.commands, command)
}

// ExecuteAll executes all health check commands concurrently
func (invoker *HealthCheckInvoker) ExecuteAll(ctx context.Context) []service.StatusLog {
	if len(invoker.commands) == 0 {
		return []service.StatusLog{}
	}

	results := make(chan service.StatusLog, len(invoker.commands))

	for _, command := range invoker.commands {
		go func(cmd HealthCheckCommand) {
			results <- cmd.Execute(ctx)
		}(command)
	}

	var statusLogs []service.StatusLog
	for i := 0; i < len(invoker.commands); i++ {
		statusLogs = append(statusLogs, <-results)
	}

	return statusLogs
}

// ExecuteSequential executes all health check commands sequentially
func (invoker *HealthCheckInvoker) ExecuteSequential(ctx context.Context) []service.StatusLog {
	var statusLogs []service.StatusLog

	for _, command := range invoker.commands {
		statusLogs = append(statusLogs, command.Execute(ctx))
	}

	return statusLogs
}
