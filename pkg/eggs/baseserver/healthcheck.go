package baseserver

import (
	"context"
	"database/sql"
	"fmt"
	"sync/atomic"
	"time"

	go_metrics "github.com/kentik/go-metrics"
)

type HealthCheckExecutor struct {
	startupDelay time.Duration
	period       time.Duration
	timeout      time.Duration
	service      Service
	latest       atomic.Value

	// metrics
	executionTimer     go_metrics.Timer
	executionMeter     go_metrics.Meter
	executionFailMeter go_metrics.Meter
	getMeter           go_metrics.Meter
	getExpiredMeter    go_metrics.Meter
}

type HealthCheckResult struct {
	Stamp   time.Time `json:"stamp"`
	Success bool      `json:"success"`
	Error   string    `json:"error"`
}

func NewHealthCheckExecutor(service Service, startupDelay time.Duration, period time.Duration, timeout time.Duration) *HealthCheckExecutor {
	registry := go_metrics.DefaultRegistry

	hce := &HealthCheckExecutor{
		startupDelay: startupDelay,
		period:       period,
		timeout:      timeout,
		service:      service,

		executionTimer:     go_metrics.NewRegisteredTimer("baseserver_healthcheck_execution_time", registry),
		executionMeter:     go_metrics.NewRegisteredMeter("baseserver_healthcheck_execution_total", registry),
		executionFailMeter: go_metrics.NewRegisteredMeter("baseserver_healthcheck_execution_fail", registry),
		getMeter:           go_metrics.NewRegisteredMeter("baseserver_healthcheck_get_total", registry),
		getExpiredMeter:    go_metrics.NewRegisteredMeter("baseserver_healthcheck_get_expired", registry),
	}
	initialResult := &HealthCheckResult{
		Stamp:   time.Now(),
		Success: true,
	}
	hce.setResult(initialResult)
	return hce
}

func (hce *HealthCheckExecutor) setResult(result *HealthCheckResult) {
	hce.latest.Store(result)
}

func (hce *HealthCheckExecutor) GetResult() *HealthCheckResult {
	hce.getMeter.Mark(1)
	latest := hce.latest.Load().(*HealthCheckResult)
	if time.Now().After(latest.Stamp.Add(hce.period * 2)) {
		hce.getExpiredMeter.Mark(1)
		replacement := &HealthCheckResult{
			Success: false,
			Stamp:   latest.Stamp,
		}
		if latest.Success {
			replacement.Error = "Latest healthcheck result was OK but has expired."
		} else {
			replacement.Error = fmt.Sprintf("Latest healthcheck result was FAILED and has expired. Error was %s", latest.Error)
		}
		return replacement
	} else {
		return latest
	}

}

func (hce *HealthCheckExecutor) Run(ctx context.Context) {
	setReady(ctx)
	time.Sleep(hce.startupDelay)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			hce.executionTimer.Time(func() {
				hce.setResult(hce.executeHealthCheck(ctx))
			})
			hce.executionMeter.Mark(1)
			time.Sleep(hce.period)
		}
	}
}

func (hce *HealthCheckExecutor) executeHealthCheck(ctx context.Context) (result *HealthCheckResult) {

	ctx, cancel := context.WithTimeout(ctx, hce.timeout)

	result = &HealthCheckResult{
		Success: true,
	}

	defer func() {
		cancel()
		if r := recover(); r != nil {
			result.Error = fmt.Sprintf("Healthcheck failed: %v", r)
			result.Success = false
		}
		result.Stamp = time.Now()
		if !result.Success {
			hce.executionFailMeter.Mark(1)
		}
	}()

	hce.service.RunHealthCheck(ctx, result)
	return
}

// Fail this healthcheck execution and panics if err is not nil.
func (hcr *HealthCheckResult) FailOnError(description string, err error) {
	if err != nil {
		panic(fmt.Errorf("%s: %v", description, err))
	}
}

// Fail this healthcheck execution with the given message and panics.
func (hcr *HealthCheckResult) Fail(message string) {
	panic(fmt.Errorf("%s", message))
}

// Fail this healthcheck execution and panics if the given sql.DB is dead
func (hcr *HealthCheckResult) FailOnSqlDbError(ctx context.Context, description string, db *sql.DB) {
	if db == nil {
		panic(fmt.Errorf("%s: nil", description))
	}

	// TODO: use PingContext and QueryContext functions here once we're sure drivers to the right thing
	// For now check context manually before queries

	hcr.FailOnError("healthcheck context closed", ctx.Err())
	hcr.FailOnError(description, db.Ping())

	// ping is optional, so let's also run a quick query
	hcr.FailOnError("healthcheck context closed", ctx.Err())
	rows, err := db.Query("SELECT NOW()")
	if rows != nil {
		hcr.FailOnError("could not close result set", rows.Close())
	}
	hcr.FailOnError(description, err)
}
