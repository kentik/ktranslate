package kt

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
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

func TestSNMPV3(t *testing.T) {
	input := []byte(`
user_name: mabel
authentication_protocol: MD5
authentication_passphrase: password123
privacy_protocol: AES
privacy_passphrase: password123
context_engine_id: aaa
context_name: ""
`)

	ms := V3SNMPConfig{}
	err := yaml.Unmarshal(input, &ms)
	assert.NoError(t, err)
	assert.Equal(t, "password123", ms.AuthenticationPassphrase)

	ser, err := yaml.Marshal(&ms)
	assert.NoError(t, err)
	assert.Equal(t, strings.TrimSpace(string(input)), strings.TrimSpace(string(ser)))

	input = []byte(`
user_name: mabel
authentication_protocol: MD5
authentication_passphrase: ${foo}
privacy_protocol: AES
privacy_passphrase: password123
context_engine_id: ${bar}
context_name: ""
`)
	t.Setenv("foo", "password123")
	t.Setenv("bar", "1234")
	err = yaml.Unmarshal(input, &ms)
	assert.NoError(t, err)
	assert.Equal(t, "password123", ms.AuthenticationPassphrase)
	assert.Equal(t, "${foo}", ms.origConf["AuthenticationPassphrase"])

	ser, err = yaml.Marshal(&ms)
	assert.NoError(t, err)
	assert.Equal(t, strings.TrimSpace(string(input)), strings.TrimSpace(string(ser)))
}
