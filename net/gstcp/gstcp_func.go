package gstcp

import (
	"crypto/rand"
	"crypto/tls"
	"github.com/jfy0o0/goStealer/errors/gserror"
	"github.com/jfy0o0/goStealer/net/gstcp/internal"
	"io"
	"net"
	"time"
)

const (
	defaultConnTimeout    = 30 * time.Second       // Default connection timeout.
	defaultRetryInternal  = 100 * time.Millisecond // Default retry interval.
	defaultReadBufferSize = 128                    // (Byte) Buffer size for reading.
)

type Retry struct {
	Count    int           // Retry count.
	Interval time.Duration // Retry interval.
}

// NewNetConn creates and returns a net.Conn with given address like "127.0.0.1:80".
// The optional parameter `timeout` specifies the timeout for dialing connection.
func NewNetConn(address string, timeout ...time.Duration) (net.Conn, error) {
	var (
		network  = `tcp`
		duration = defaultConnTimeout
	)
	if len(timeout) > 0 {
		duration = timeout[0]
	}
	conn, err := net.DialTimeout(network, address, duration)
	if err != nil {
		err = gserror.Wrapf(
			err,
			`net.DialTimeout failed with network "%s", address "%s", timeout "%s"`,
			network, address, duration,
		)
	}
	return conn, err
}

// NewNetConnTLS creates and returns a TLS net.Conn with given address like "127.0.0.1:80".
// The optional parameter `timeout` specifies the timeout for dialing connection.
func NewNetConnTLS(address string, tlsConfig *tls.Config, timeout ...time.Duration) (net.Conn, error) {
	var (
		network = `tcp`
		dialer  = &net.Dialer{
			Timeout: defaultConnTimeout,
		}
	)
	if len(timeout) > 0 {
		dialer.Timeout = timeout[0]
	}
	conn, err := tls.DialWithDialer(dialer, network, address, tlsConfig)
	if err != nil {
		err = gserror.Wrapf(
			err,
			`tls.DialWithDialer failed with network "%s", address "%s", timeout "%s", tlsConfig "%v"`,
			network, address, dialer.Timeout, tlsConfig,
		)
	}
	return conn, err
}

// NewNetConnKeyCrt creates and returns a TLS net.Conn with given TLS certificate and key files
// and address like "127.0.0.1:80". The optional parameter `timeout` specifies the timeout for
// dialing connection.
func NewNetConnKeyCrt(addr, crtFile, keyFile string, timeout ...time.Duration) (net.Conn, error) {
	tlsConfig, err := LoadKeyCrt(crtFile, keyFile)
	if err != nil {
		return nil, err
	}
	return NewNetConnTLS(addr, tlsConfig, timeout...)
}

// Send creates connection to `address`, writes `data` to the connection and then closes the connection.
// The optional parameter `retry` specifies the retry policy when fails in writing data.
func Send(address string, data []byte, retry ...Retry) error {
	conn, err := NewConn(address)
	if err != nil {
		return err
	}
	defer conn.rawConn.Close()
	return conn.Send(data, retry...)
}

// SendRecv creates connection to `address`, writes `data` to the connection, receives response
// and then closes the connection.
//
// The parameter `length` specifies the bytes count waiting to receive. It receives all buffer content
// and returns if `length` is -1.
//
// The optional parameter `retry` specifies the retry policy when fails in writing data.
func SendRecv(address string, data []byte, length int, retry ...Retry) ([]byte, error) {
	conn, err := NewConn(address)
	if err != nil {
		return nil, err
	}
	defer conn.rawConn.Close()
	return conn.SendRecv(data, length, retry...)
}

// SendWithTimeout does Send logic with writing timeout limitation.
func SendWithTimeout(address string, data []byte, timeout time.Duration, retry ...Retry) error {
	conn, err := NewConn(address)
	if err != nil {
		return err
	}
	defer conn.rawConn.Close()
	return conn.SendWithTimeout(data, timeout, retry...)
}

// SendRecvWithTimeout does SendRecv logic with reading timeout limitation.
func SendRecvWithTimeout(address string, data []byte, receive int, timeout time.Duration, retry ...Retry) ([]byte, error) {
	conn, err := NewConn(address)
	if err != nil {
		return nil, err
	}
	defer conn.rawConn.Close()
	return conn.SendRecvWithTimeout(data, receive, timeout, retry...)
}

// isTimeout checks whether given `err` is a timeout error.
func isTimeout(err error) bool {
	if err == nil {
		return false
	}
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return true
	}
	return false
}

// LoadKeyCrt creates and returns a TLS configuration object with given certificate and key files.
func LoadKeyCrt(crtFile, keyFile string) (*tls.Config, error) {
	crt, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		return nil, gserror.Wrapf(err,
			`tls.LoadX509KeyPair failed for certFile "%s" and keyFile "%s"`,
			crtFile, keyFile,
		)
	}
	tlsConfig := &tls.Config{}
	tlsConfig.Certificates = []tls.Certificate{crt}
	tlsConfig.Time = time.Now
	tlsConfig.Rand = rand.Reader
	return tlsConfig, nil
}

// MustGetFreePort performs as GetFreePort, but it panics is any error occurs.
func MustGetFreePort() int {
	port, err := GetFreePort()
	if err != nil {
		panic(err)
	}
	return port
}

// GetFreePort retrieves and returns a port that is free.
func GetFreePort() (port int, err error) {
	var (
		network = `tcp`
		address = `:0`
	)
	resolvedAddr, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		return 0, gserror.Wrapf(
			err,
			`net.ResolveTCPAddr failed for network "%s", address "%s"`,
			network, address,
		)
	}
	l, err := net.ListenTCP(network, resolvedAddr)
	if err != nil {
		return 0, gserror.Wrapf(
			err,
			`net.ListenTCP failed for network "%s", address "%s"`,
			network, resolvedAddr.String(),
		)
	}
	port = l.Addr().(*net.TCPAddr).Port
	err = l.Close()
	return
}

// GetFreePorts retrieves and returns specified number of ports that are free.
func GetFreePorts(count int) (ports []int, err error) {
	var (
		network = `tcp`
		address = `:0`
	)
	for i := 0; i < count; i++ {
		resolvedAddr, err := net.ResolveTCPAddr(network, address)
		if err != nil {
			return nil, gserror.Wrapf(
				err,
				`net.ResolveTCPAddr failed for network "%s", address "%s"`,
				network, address,
			)
		}
		l, err := net.ListenTCP(network, resolvedAddr)
		if err != nil {
			return nil, gserror.Wrapf(
				err,
				`net.ListenTCP failed for network "%s", address "%s"`,
				network, resolvedAddr.String(),
			)
		}
		ports = append(ports, l.Addr().(*net.TCPAddr).Port)
		_ = l.Close()
	}
	return ports, nil
}

func UpgradeConnAsClient(conn net.Conn, cliTlsConfig ...*tls.Config) *Conn {
	config := internal.GetClientTlsConfig()
	if len(cliTlsConfig) > 0 {
		config = cliTlsConfig[0]
	}
	return NewConnByNetConn(tls.Client(conn, config))
}

func UpgradeConnAsServer(conn net.Conn, svrTlsConfig ...*tls.Config) *Conn {
	config := internal.GetServerTlsConfig()
	if len(svrTlsConfig) > 0 {
		config = svrTlsConfig[0]
	}
	return NewConnByNetConn(tls.Server(conn, config))
}

// RelayConnection relay copies between left and right bidirectionally. Returns number of
// bytes copied from right to left, from left to right, and any error occurred.
func RelayConnection(left, right net.Conn) (int64, int64, error) {
	type res struct {
		N   int64
		Err error
	}
	ch := make(chan res)

	go func() {
		n, err := io.Copy(right, left)
		_ = right.SetDeadline(time.Now()) // wake up the other goroutine blocking on right
		_ = left.SetDeadline(time.Now())  // wake up the other goroutine blocking on left
		ch <- res{n, err}
	}()

	n, err := io.Copy(left, right)
	_ = right.SetDeadline(time.Now()) // wake up the other goroutine blocking on right
	_ = left.SetDeadline(time.Now())  // wake up the other goroutine blocking on left
	rs := <-ch

	if err == nil {
		err = rs.Err
	}
	return n, rs.N, err
}
