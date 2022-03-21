package kt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsPollReady(t *testing.T) {
	// Empty mib always returns true.
	mib := &Mib{}
	assert.True(t, mib.IsPollReady())
	assert.True(t, mib.IsPollReady())
	assert.True(t, mib.IsPollReady())

	// Now, set a poll duration.
	mib.PollDur = time.Duration(10) * time.Second
	assert.True(t, mib.IsPollReady())  // first poll is good.
	assert.False(t, mib.IsPollReady()) // Skip the 2nd.
	assert.False(t, mib.IsPollReady()) // Skip the 2nd.
}

func TestGetName(t *testing.T) {
	mib := &Mib{
		Tag:  "foo",
		Name: "name",
	}
	assert.Equal(t, "foo", mib.GetName())
	mib = nil
	assert.Equal(t, "missing_mib", mib.GetName())
	mib = &Mib{
		Name: "bar",
	}
	assert.Equal(t, "bar", mib.GetName())
}
