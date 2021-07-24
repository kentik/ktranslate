package reuse

import (
	"context"
	"net"
	"syscall"
	"time"
)

var (
	Config           = net.ListenConfig{Control: onConnSetup}
	ReusableListener = func(network, host string) (net.Listener, error) {
		return Config.Listen(context.Background(), network, host)
	}
)

func onConnSetup(_ string, _ string, c syscall.RawConn) (err error) {
	return
}

// Copied from net/http/server.go
//
// TcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type TcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln TcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
