package baseserver

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func doGet(t *testing.T, addr net.Addr, path string) (int, string) {
	url := fmt.Sprintf("http://%s/%s", addr, path)
	t.Logf("doGet: hitting %s", url)
	resp, err := http.Get(url)
	assert.Nil(t, err, "doGet: http.Get failed")
	defer resp.Body.Close()
	payload, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err, "doGet: io.ReadAll failed")
	return resp.StatusCode, string(payload)
}

func TestMetaServer(t *testing.T) {
	BaseServerConfigurationDefaults.SkipParseFlags = true
	BaseServerConfigurationDefaults.LogToStdout = false
	BaseServerConfigurationDefaults.ShutdownSettleTime = 1 * time.Millisecond

	BaseServerConfigurationDefaults.HealthCheckStartupDelay = 50 * time.Millisecond
	BaseServerConfigurationDefaults.HealthCheckPeriod = 1 * time.Second
	BaseServerConfigurationDefaults.HealthCheckTimeout = 1 * time.Second

	bs := Boilerplate("dumdum", versionInfo,nil)
	assert.NotNil(t, bs, "Boilerplate should return a non-nil value")

	dumdum := NewDummyService(t, false, 1*time.Millisecond)

	go func() {
		bs.Run(dumdum)
	}()

	fmt.Println("waiting until ready")
	bs.WaitUntilReady(time.Second)

	assert.NotNil(t, bs.metaServer.listenAddr, "metaserver listen addr must not be nil")

	status, _ := doGet(t, bs.metaServer.listenAddr, "healthcheck")
	assert.Equal(t, 200, status, "first healthcheck should succeed")

	time.Sleep(75 * time.Millisecond)
	status, _ = doGet(t, bs.metaServer.listenAddr, "healthcheck")
	assert.Equal(t, 500, status, "second healthcheck should fail")

	status, _ = doGet(t, bs.metaServer.listenAddr, "metrics")
	assert.Equal(t, 200, status, "metrics endpoint should succeed")

	status, _ = doGet(t, bs.metaServer.listenAddr, "version")
	assert.Equal(t, 200, status, "version endpoint should succeed")

	status, _ = doGet(t, bs.metaServer.listenAddr, "env")
	assert.Equal(t, 200, status, "env endpoint should succeed")

	resetGlobalBaseServer()
}
