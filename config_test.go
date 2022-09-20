package ktranslate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestKtranslateLoadConfigDefault(t *testing.T) {
	conf := `listenaddr: 127.0.0.1:8081
dns: "1.1.1.1"
processingthreads: 2
inputthreads: 2
maxthreads: 2
format: flat_json
formatrollup: ""
compression: none
maxflowspermessage: 10000
rollupinterval: 0
rollupandalpha: false
sinks: "http"
samplerate: 1
samplemin: 1
enablesnmpdiscovery: false
kentikemail: ""
kentikapitoken: ""
kentikplan: 0
apibaseurl: https://api.kentik.com
enableteelogs: false
enablehttpinput: true
server:
    servicename: "ktranslate"
    loglevel: info
    logtostdout: true
    metricsendpoint: none
    metalistenaddr: ""
    ollydataset: ""
    ollywritekey: ""
httpsink:
  target: "http://127.0.0.1:8080"
  headers: ["TEST1:foo", "TEST2:bar"]
`
	var cfg *Config
	if err := yaml.Unmarshal([]byte(conf), &cfg); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, cfg.ListenAddr, "127.0.0.1:8081")
	assert.Equal(t, cfg.ProcessingThreads, 2)
}

func TestKtranslateLoadConfigMultiSink(t *testing.T) {
	conf := `listenaddr: 127.0.0.1:8081
dns: "1.1.1.1"
processingthreads: 2
inputthreads: 2
maxthreads: 2
format: flat_json
formatrollup: ""
compression: none
maxflowspermessage: 10000
rollupinterval: 0
rollupandalpha: false
sinks: "http"
samplerate: 1
samplemin: 1
enablesnmpdiscovery: false
kentikemail: ""
kentikapitoken: ""
kentikplan: 0
apibaseurl: https://api.kentik.com
enableteelogs: false
enablehttpinput: true
server:
    servicename: "ktranslate"
    loglevel: info
    logtostdout: true
    metricsendpoint: none
    metalistenaddr: ""
    ollydataset: ""
    ollywritekey: ""
multisink:
 filesinks:
   - path: "/foo"
   - path: "/bar"
 httpsinks:
   - target: "127.0.0.1:10000"
   - target: "127.0.0.1:10001"

`
	var cfg *Config
	if err := yaml.Unmarshal([]byte(conf), &cfg); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, cfg.ListenAddr, "127.0.0.1:8081")
	assert.Equal(t, cfg.ProcessingThreads, 2)
	assert.NotNil(t, cfg.MultiSink)
	assert.NotNil(t, cfg.MultiSink.FileSinks)
	assert.Equal(t, len(cfg.MultiSink.FileSinks), 2)
	assert.NotNil(t, cfg.MultiSink.HTTPSinks)
	assert.Equal(t, len(cfg.MultiSink.HTTPSinks), 2)
}
