package healthcheck

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"runtime"
	"time"

	"github.com/kentik/ktranslate/pkg/util/logger"
	"github.com/kentik/ktranslate/pkg/util/reuse"
)

const (
	LOG_PREFIX      = "[HealthCheck] "
	PPROF_BIND_ADDR = "PPROF_BIND_ADDR"
	INIT_TIME       = 10 * time.Second
)

var (
	GOOD = []byte("GOOD\n")
	BAD  = []byte("BAD\n")
)

func nilStatus() []byte {
	return GOOD
}

func nilCmd(cmd []byte) []byte {
	return []byte(fmt.Sprintf("Unknown command: %s\n", string(cmd)))
}

// GetMemStats returns a simple message describing the state of heap memory, as seen by the runtime.
// All values are in MB.
// - Sys:          bytes obtained from system
// - HeapSys:      bytes obtained from system
// - HeapAlloc:    bytes allocated and not yet freed
// - HeapIdle:     bytes in idle spans
// - HeapReleased: bytes released to the OS
func GetMemStats() string {
	mb := uint64(1000000)
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	return fmt.Sprintf("Sys: %d, HeapSys: %d, HeapAlloc: %d, HeapIdle: %d, HeapReleased: %d",
		memStats.Sys/mb, memStats.HeapSys/mb, memStats.HeapAlloc/mb, memStats.HeapIdle/mb, memStats.HeapReleased/mb)
}

func peekCmd(c net.Conn) ([]byte, error) {
	c.SetReadDeadline(time.Now().Add(time.Millisecond * 500))
	r := bufio.NewReaderSize(c, 200) // maximum command length == 200

	if _, err := r.Peek(1); err == nil {
		// Extend read deadline
		c.SetReadDeadline(time.Now().Add(time.Second * 5))
		if buf, _, err := r.ReadLine(); err != nil {
			return nil, err
		} else {
			return buf, nil
		}
	}
	return nil, nil
}

func Run(host string, statusReport func() []byte, handleCmd func([]byte) []byte, log *logger.Logger) {
	RunWithListener(net.Listen, host, statusReport, handleCmd, log)
}

func RunWithListener(listenFunc func(network, host string) (net.Listener, error),
	host string, statusReport func() []byte, handleCmd func([]byte) []byte, log *logger.Logger) {

	// Sleep for the time needed to settle things.
	time.Sleep(INIT_TIME)

	l, err := listenFunc("tcp", host)
	if err != nil {
		log.Error(LOG_PREFIX, "Error Binding to %s: %v", host, err)
		return
	}

	if statusReport == nil {
		statusReport = nilStatus
	}

	if handleCmd == nil {
		handleCmd = nilCmd
	}

	// Start up a pprof server as well.
	pprofBind := os.Getenv(PPROF_BIND_ADDR)
	if pprofBind != "" {
		go func() {
			mux := http.NewServeMux()
			mux.HandleFunc("/debug/pprof/", pprof.Index)
			mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
			mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
			mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
			mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

			ln, err := listenFunc("tcp", pprofBind)
			if err != nil {
				log.Error(LOG_PREFIX, "Healthcheck pprof server listener error: %s", err)
				return
			}

			log.Info(LOG_PREFIX, "pprof listening at %s", pprofBind)
			srv := &http.Server{Addr: pprofBind, Handler: mux}
			log.Error(LOG_PREFIX, "%s", srv.Serve(reuse.TcpKeepAliveListener{ln.(*net.TCPListener)}))
		}()
	}

	log.Info(LOG_PREFIX, "HC online at %s", host)
	defer l.Close()
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Error(LOG_PREFIX, "Error Accepting request: %v", err)
			continue
		}

		cmd, err := peekCmd(conn)
		if err != nil {
			log.Error(LOG_PREFIX, "Error reading command: %s", err)
			conn.Close()
			continue
		}

		if cmd != nil && !bytes.Equal(cmd, []byte("status")) {
			resp := handleCmd(cmd)
			conn.Write(resp)
		} else {
			resp := statusReport()
			conn.Write(resp)
		}

		conn.Close()
	}
}
