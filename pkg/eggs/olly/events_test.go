package olly

import (
	"testing"

	"github.com/kentik/ktranslate/pkg/eggs/timing"
	"github.com/stretchr/testify/assert"

	"github.com/honeycombio/libhoney-go"
)

func TestData2Map(t *testing.T) {
	m := data2map("lat", 1.1, "lon", 2.2, "name", "nowhere")
	assert.Equal(t, 3, len(m))
	assert.Equal(t, 1.1, m["lat"])
	assert.Equal(t, 2.2, m["lon"])
	assert.Equal(t, "nowhere", m["name"])

}

func TestClosingTimeout(t *testing.T) {
	// send to a bad destination so events get buffered
	Init("test", "1.0", libhoney.Config{
		APIHost:  "http://api.kentik.com:1234",
		WriteKey: "lol",
		Dataset:  "yeah",
	})
	b := NewBuilder()
	for i := 0; i < 10000; i++ {
		Send(newEvent(b, Op("test.test")))
	}

	c := timing.StartChrono()
	Close() // close fill try to flush events but should give up after CloseFlushTimeout
	assert.True(t, c.FinishDur() < 5*CloseFlushTimeout)
}
