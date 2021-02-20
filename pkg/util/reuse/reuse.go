package reuse

import (
	"context"
	"net"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

var (
	Config           = net.ListenConfig{Control: onConnSetup}
	ReusableListener = func(network, host string) (net.Listener, error) {
		return Config.Listen(context.Background(), network, host)
	}
)

func onConnSetup(_ string, _ string, c syscall.RawConn) (err error) {
	c.Control(func(fd uintptr) {
		err = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, unix.SO_REUSEPORT, 1)
		// ENOPROTOOPT means that SO_REUSEPORT is not implemented in the current
		// kernel.  All we can do is shrug and move on.  Should only happen in
		// chfagent (which doesn't do handoffs), on fairly old (pre-2013) Linux
		// kernels.
		if errno, ok := err.(syscall.Errno); ok && errno == syscall.ENOPROTOOPT {
			err = nil
		}
	})
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
