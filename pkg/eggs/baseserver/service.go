package baseserver

import (
	"context"
	"net/http"

	"github.com/judwhite/go-svc"
)

// Service interface - baseserver/metaserver use these methods to interact with actual services
type Service interface {
	// Run this service until the given context is closed, after which the service cannot be re-started. Blocks.
	Run(ctx context.Context) error

	GetStatus() []byte                                             // for legacy healthcheck
	RunHealthCheck(ctx context.Context, result *HealthCheckResult) // new style healthcheck

	HttpInfo(w http.ResponseWriter, req *http.Request) // Hook to provide http info via metaserver

	// These are needed in case we are running under windows.
	Init(svc.Environment) error

	Start() error

	Stop() error
}
