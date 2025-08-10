package checker

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
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
			Error:       fmt.Errorf("failed to create request: %w", err).Error(),
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
		statusLog.Error = fmt.Errorf("request failed after %d attempts: %w", maxRetries, lastErr).Error()
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

// WorkerPoolConfig holds configuration for the worker pool
type WorkerPoolConfig struct {
	WorkerCount        int           // Number of workers in the pool
	MaxConcurrent      int           // Maximum concurrent checks per service
	GlobalTimeout      time.Duration // Global timeout for all checks
	PerProbeTimeout    time.Duration // Individual probe timeout
	JitterMaxDuration  time.Duration // Maximum jitter to add
	RetryAttempts      int           // Number of retry attempts
	RetryBackoffFactor float64       // Backoff multiplier for retries
	RetryInitialDelay  time.Duration // Initial retry delay
}

// DefaultWorkerPoolConfig returns default configuration
func DefaultWorkerPoolConfig() WorkerPoolConfig {
	return WorkerPoolConfig{
		WorkerCount:        10,
		MaxConcurrent:      5,
		GlobalTimeout:      5 * time.Minute,
		PerProbeTimeout:    30 * time.Second,
		JitterMaxDuration:  time.Second,
		RetryAttempts:      3,
		RetryBackoffFactor: 2.0,
		RetryInitialDelay:  500 * time.Millisecond,
	}
}

// HealthCheckInvoker manages health check commands with bounded worker pool
type HealthCheckInvoker struct {
	commands []HealthCheckCommand
	config   WorkerPoolConfig
}

// NewHealthCheckInvoker creates a new health check invoker with default config
func NewHealthCheckInvoker() *HealthCheckInvoker {
	return NewHealthCheckInvokerWithConfig(DefaultWorkerPoolConfig())
}

// NewHealthCheckInvokerWithConfig creates a new health check invoker with custom config
func NewHealthCheckInvokerWithConfig(config WorkerPoolConfig) *HealthCheckInvoker {
	return &HealthCheckInvoker{
		commands: make([]HealthCheckCommand, 0),
		config:   config,
	}
}

// AddCommand adds a health check command
func (invoker *HealthCheckInvoker) AddCommand(command HealthCheckCommand) {
	invoker.commands = append(invoker.commands, command)
}

// ExecuteAll executes all health check commands using bounded worker pool
func (invoker *HealthCheckInvoker) ExecuteAll(ctx context.Context) []service.StatusLog {
	if len(invoker.commands) == 0 {
		return []service.StatusLog{}
	}

	// Create context with global timeout
	globalCtx, cancel := context.WithTimeout(ctx, invoker.config.GlobalTimeout)
	defer cancel()

	return invoker.executeWithWorkerPool(globalCtx)
}

// executeWithWorkerPool implements bounded worker pool execution
func (invoker *HealthCheckInvoker) executeWithWorkerPool(ctx context.Context) []service.StatusLog {
	jobs := make(chan HealthCheckCommand, len(invoker.commands))
	results := make(chan service.StatusLog, len(invoker.commands))
	
	// Create semaphore for per-service concurrency control
	semaphore := make(chan struct{}, invoker.config.MaxConcurrent)

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < invoker.config.WorkerCount; i++ {
		wg.Add(1)
		go invoker.worker(ctx, jobs, results, semaphore, &wg)
	}

	// Send jobs
	go func() {
		defer close(jobs)
		for _, cmd := range invoker.commands {
			select {
			case jobs <- cmd:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Wait for all workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var statusLogs []service.StatusLog
	for result := range results {
		statusLogs = append(statusLogs, result)
	}

	return statusLogs
}

// worker processes health check commands with jitter, timeout, and retry logic
func (invoker *HealthCheckInvoker) worker(ctx context.Context, jobs <-chan HealthCheckCommand, results chan<- service.StatusLog, semaphore chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case cmd, ok := <-jobs:
			if !ok {
				return
			}

			// Acquire semaphore for per-service concurrency control
			select {
			case semaphore <- struct{}{}:
			case <-ctx.Done():
				return
			}

			// Add jitter before executing
			if invoker.config.JitterMaxDuration > 0 {
				jitter := time.Duration(rand.Int63n(int64(invoker.config.JitterMaxDuration)))
				time.Sleep(jitter)
			}

			// Execute with per-probe timeout and retry logic
			result := invoker.executeCommandWithRetry(ctx, cmd)
			
			// Release semaphore
			<-semaphore

			select {
			case results <- result:
			case <-ctx.Done():
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

// executeCommandWithRetry executes a command with retry and backoff logic
func (invoker *HealthCheckInvoker) executeCommandWithRetry(ctx context.Context, cmd HealthCheckCommand) service.StatusLog {
	var lastResult service.StatusLog
	delay := invoker.config.RetryInitialDelay

	for attempt := 0; attempt < invoker.config.RetryAttempts; attempt++ {
		// Create context with per-probe timeout
		probeCtx, cancel := context.WithTimeout(ctx, invoker.config.PerProbeTimeout)
		
		// Execute command
		result := cmd.Execute(probeCtx)
		cancel()

		// If successful or context cancelled, return immediately
		if result.Status == statusOperational || ctx.Err() != nil {
			return result
		}

		lastResult = result

		// If not the last attempt, wait before retrying
		if attempt < invoker.config.RetryAttempts-1 {
			select {
			case <-time.After(delay):
				// Increase delay with backoff factor
				delay = time.Duration(float64(delay) * invoker.config.RetryBackoffFactor)
			case <-ctx.Done():
				return lastResult
			}
		}
	}

	return lastResult
}

// ExecuteSequential executes all health check commands sequentially
func (invoker *HealthCheckInvoker) ExecuteSequential(ctx context.Context) []service.StatusLog {
	var statusLogs []service.StatusLog

	for _, command := range invoker.commands {
		statusLogs = append(statusLogs, command.Execute(ctx))
	}

	return statusLogs
}
