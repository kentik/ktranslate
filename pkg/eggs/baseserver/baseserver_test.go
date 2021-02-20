package baseserver

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/version"
	"github.com/stretchr/testify/assert"
)

var versionInfo = version.VersionInfo{}

type DummyService struct {
	t                *testing.T
	alwaysUp         bool
	healthcheckSleep time.Duration
	contextGotDone   bool
}

func NewDummyService(t *testing.T, alwaysUp bool, healthcheckSleep time.Duration) *DummyService {
	return &DummyService{
		t:                t,
		alwaysUp:         alwaysUp,
		healthcheckSleep: healthcheckSleep,
		contextGotDone:   false,
	}
}

func (dumdum *DummyService) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			dumdum.contextGotDone = true
			dumdum.t.Log("DummyService.Run() returning")
			return nil
		}
	}
}

func (dumdum *DummyService) GetStatus() []byte {
	return nil
}

func (dumdum *DummyService) RunHealthCheck(ctx context.Context, result *HealthCheckResult) {
	dumdum.t.Log("DummyService.RunHealthCheck entering")
	time.Sleep(dumdum.healthcheckSleep)
	if !dumdum.alwaysUp {
		result.Fail("This DummyService was configured to be always down")
	}
	dumdum.t.Log("DummyService.RunHealthCheck returning")
}

func (dumdum *DummyService) Close() {
	dumdum.t.Log("DummyService.Close returning")
}

func (dumdum *DummyService) HttpInfo(w http.ResponseWriter, req *http.Request) {
	// Noop
	dumdum.t.Log("DummyService.HttpInfo returning")
}

func TestBoilerplateAndShutdown(t *testing.T) {
	BaseServerConfigurationDefaults.SkipParseFlags = true
	BaseServerConfigurationDefaults.LogToStdout = false
	BaseServerConfigurationDefaults.ShutdownSettleTime = 1 * time.Millisecond
	bs := Boilerplate("dumdum", versionInfo, nil)
	assert.NotNil(t, bs, "Boilerplate should return a non-nil value")

	assert.NotNil(t, GetGlobalBaseServer())

	dumdum := NewDummyService(t, true, 10*time.Millisecond)

	// run bs.Run in the background, sleep a bit then shutdown
	done := make(chan struct{})
	go func() {
		bs.Run(dumdum)
		close(done)
	}()

	time.Sleep(20 * time.Millisecond)
	bs.Shutdown("test")
	<-done
	assert.True(t, dumdum.contextGotDone, "service context should get closed during shutdown")
	resetGlobalBaseServer()
}

func TestSignalShutdown(t *testing.T) {
	BaseServerConfigurationDefaults.SkipParseFlags = true
	BaseServerConfigurationDefaults.LogToStdout = false
	BaseServerConfigurationDefaults.ShutdownSettleTime = 1 * time.Millisecond
	bs := Boilerplate("dumdum", versionInfo, nil)
	assert.NotNil(t, bs, "Boilerplate returns a non-nil value")

	dumdum := NewDummyService(t, true, 10*time.Millisecond)

	// run bs.Run in the background, sleep a bit then shutdown
	go func() {
		bs.Run(dumdum)
	}()

	time.Sleep(20 * time.Millisecond)

	process, err := os.FindProcess(os.Getpid())
	assert.Nil(t, err, "os.FindProcess returns a non-nil value")
	process.Signal(os.Interrupt) // nolint: errcheck

	time.Sleep(1 * time.Millisecond)
	assert.True(t, dumdum.contextGotDone, "service context should get closed during shutdown")
	resetGlobalBaseServer()
}
