package flock

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFlock(t *testing.T) {
	defer os.Remove("TestFlock.lock")

	lock, err := New("TestFlock.lock")
	assert.NotNil(t, lock)
	assert.NoError(t, err)
	defer lock.Close()

	err = lock.LockWithTimeout(time.Millisecond, time.Millisecond)
	if assert.NoError(t, err) {
		lock.Unlock()
	} else {
		assert.FailNow(t, "could not get lock")
	}
}
