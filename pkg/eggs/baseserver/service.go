package baseserver

import (
	"context"
	"net/http"
)

// Service interface - baseserver/metaserver use these methods to interact with actual services
type Service interface {
	// Run this service until the given context is closed, after which the service cannot be re-started. Blocks.
	Run(ctx context.Context) error

	GetStatus() []byte                                             // for legacy healthcheck
	RunHealthCheck(ctx context.Context, result *HealthCheckResult) // new style healthcheck

	Close() // deprecated; select fron ctx.Done() in Run() instead

	HttpInfo(w http.ResponseWriter, req *http.Request) // Hook to provide http info via metaserver
}
