package baseserver

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckExecutor(t *testing.T) {

	ctx := context.Background()

	dumdum := NewDummyService(t, false, 10*time.Millisecond)
	hce := NewHealthCheckExecutor(dumdum, 5*time.Millisecond, 5*time.Millisecond, 5*time.Second)
	go hce.Run(ctx)

	assert.True(t, hce.GetResult().Success, "initial result should be successful")
	time.Sleep(10 * time.Millisecond)
	assert.False(t, hce.GetResult().Success, "subsequent results should not be successful")
}

func TestHealthCheckExecutor_Expiration(t *testing.T) {

	ctx := context.Background()

	dumdum := NewDummyService(t, true, 30*time.Millisecond)

	hce := NewHealthCheckExecutor(dumdum, 5*time.Millisecond, 5*time.Millisecond, 5*time.Second)
	go hce.Run(ctx)

	assert.True(t, hce.GetResult().Success, "initial result should be successful")
	time.Sleep(10 * time.Millisecond)
	assert.False(t, hce.GetResult().Success, "second result should not be successful, because of expiration")
	assert.Equal(t, hce.GetResult().Error, "Latest healthcheck result was OK but has expired.", "Error should be about expiration")

	seenSuccess := false
	for i := 0; i <= 20; i++ {
		time.Sleep(5 * time.Millisecond)
		if hce.GetResult().Success {
			seenSuccess = true
			break
		}
	}
	assert.True(t, seenSuccess, "every so often we should hit a non-expired, successful result")
}
