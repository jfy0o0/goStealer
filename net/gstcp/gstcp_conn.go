package gstcp

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"github.com/jfy0o0/goStealer/errors/gserror"
	"io"
	"net"
	"time"
)

// Conn is the TCP connection object.
type Conn struct {
	rawConn           net.Conn      // Underlying TCP connection object.
	reader            *bufio.Reader // Buffer reader for connection.
	receiveDeadline   time.Time     // Timeout point for reading.
	sendDeadline      time.Time     // Timeout point for writing.
	receiveBufferWait time.Duration // Interval duration for reading buffer.
}

const (
	// Default interval for reading buffer.
	receiveAllWaitTimeout = time.Millisecond
)

// NewConn creates and returns a new connection with given address.
func NewConn(addr string, timeout ...time.Duration) (*Conn, error) {
	if conn, err := NewNetConn(addr, timeout...); err == nil {
		return NewConnByNetConn(conn), nil
	} else {
		return nil, err
	}
}

// NewConnTLS creates and returns a new TLS connection
// with given address and TLS configuration.
func NewConnTLS(addr string, tlsConfig *tls.Config) (*Conn, error) {
	if conn, err := NewNetConnTLS(addr, tlsConfig); err == nil {
		return NewConnByNetConn(conn), nil
	} else {
		return nil, err
	}
}

// NewConnKeyCrt creates and returns a new TLS connection
// with given address and TLS certificate and key files.
func NewConnKeyCrt(addr, crtFile, keyFile string) (*Conn, error) {
	if conn, err := NewNetConnKeyCrt(addr, crtFile, keyFile); err == nil {
		return NewConnByNetConn(conn), nil
	} else {
		return nil, err
	}
}

// NewConnByNetConn creates and returns a TCP connection object with given net.Conn object.
func NewConnByNetConn(conn net.Conn) *Conn {
	return &Conn{
		rawConn:           conn,
		reader:            bufio.NewReader(conn),
		receiveDeadline:   time.Time{},
		sendDeadline:      time.Time{},
		receiveBufferWait: receiveAllWaitTimeout,
	}
}

// Send writes data to remote address.
func (c *Conn) Send(data []byte, retry ...Retry) error {
	for {
		if _, err := c.rawConn.Write(data); err != nil {
			// Connection closed.
			if err == io.EOF {
				return err
			}
			// Still failed even after retrying.
			if len(retry) == 0 || retry[0].Count == 0 {
				err = gserror.Wrap(err, `Write data failed`)
				return err
			}
			if len(retry) > 0 {
				retry[0].Count--
				if retry[0].Interval == 0 {
					retry[0].Interval = defaultRetryInternal
				}
				time.Sleep(retry[0].Interval)
			}
		} else {
			return nil
		}
	}
}

// Recv receives and returns data from the connection.
//
// Note that,
//  1. If length = 0, which means it receives the data from current buffer and returns immediately.
//  2. If length < 0, which means it receives all data from connection and returns it until no data
//     from connection. Developers should notice the package parsing yourself if you decide receiving
//     all data from buffer.
//  3. If length > 0, which means it blocks reading data from connection until length size was received.
//     It is the most commonly used length value for data receiving.
func (c *Conn) Recv(length int, retry ...Retry) ([]byte, error) {
	var (
		err        error  // Reading error.
		size       int    // Reading size.
		index      int    // Received size.
		buffer     []byte // Buffer object.
		bufferWait bool   // Whether buffer reading timeout set.
	)
	if length > 0 {
		buffer = make([]byte, length)
	} else {
		buffer = make([]byte, defaultReadBufferSize)
	}

	for {
		if length < 0 && index > 0 {
			bufferWait = true
			if err = c.rawConn.SetReadDeadline(time.Now().Add(c.receiveBufferWait)); err != nil {
				err = gserror.Wrap(err, `SetReadDeadline for connection failed`)
				return nil, err
			}
		}
		size, err = c.reader.Read(buffer[index:])
		if size > 0 {
			index += size
			if length > 0 {
				// It reads til `length` size if `length` is specified.
				if index == length {
					break
				}
			} else {
				if index >= defaultReadBufferSize {
					// If it exceeds the buffer size, it then automatically increases its buffer size.
					buffer = append(buffer, make([]byte, defaultReadBufferSize)...)
				} else {
					// It returns immediately if received size is lesser than buffer size.
					if !bufferWait {
						break
					}
				}
			}
		}
		if err != nil {
			// Connection closed.
			if err == io.EOF {
				break
			}
			// Re-set the timeout when reading data.
			if bufferWait && isTimeout(err) {
				if err = c.rawConn.SetReadDeadline(c.receiveDeadline); err != nil {
					err = gserror.Wrap(err, `SetReadDeadline for connection failed`)
					return nil, err
				}
				err = nil
				break
			}
			if len(retry) > 0 {
				// It fails even it retried.
				if retry[0].Count == 0 {
					break
				}
				retry[0].Count--
				if retry[0].Interval == 0 {
					retry[0].Interval = defaultRetryInternal
				}
				time.Sleep(retry[0].Interval)
				continue
			}
			break
		}
		// Just read once from buffer.
		if length == 0 {
			break
		}
	}
	return buffer[:index], err
}

// RecvLine reads data from the connection until reads char '\n'.
// Note that the returned result does not contain the last char '\n'.
func (c *Conn) RecvLine(retry ...Retry) ([]byte, error) {
	var (
		err    error
		buffer []byte
		data   = make([]byte, 0)
	)
	for {
		buffer, err = c.Recv(1, retry...)
		if len(buffer) > 0 {
			if buffer[0] == '\n' {
				data = append(data, buffer[:len(buffer)-1]...)
				break
			} else {
				data = append(data, buffer...)
			}
		}
		if err != nil {
			break
		}
	}
	return data, err
}

// RecvTill reads data from the connection until reads bytes `til`.
// Note that the returned result contains the last bytes `til`.
func (c *Conn) RecvTill(til []byte, retry ...Retry) ([]byte, error) {
	var (
		err    error
		buffer []byte
		data   = make([]byte, 0)
		length = len(til)
	)
	for {
		buffer, err = c.Recv(1, retry...)
		if len(buffer) > 0 {
			if length > 0 &&
				len(data) >= length-1 &&
				buffer[0] == til[length-1] &&
				bytes.EqualFold(data[len(data)-length+1:], til[:length-1]) {
				data = append(data, buffer...)
				break
			} else {
				data = append(data, buffer...)
			}
		}
		if err != nil {
			break
		}
	}
	return data, err
}

// RecvWithTimeout reads data from the connection with timeout.
func (c *Conn) RecvWithTimeout(length int, timeout time.Duration, retry ...Retry) (data []byte, err error) {
	if err = c.SetReceiveDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}
	defer c.SetReceiveDeadline(time.Time{})
	data, err = c.Recv(length, retry...)
	return
}

// SendWithTimeout writes data to the connection with timeout.
func (c *Conn) SendWithTimeout(data []byte, timeout time.Duration, retry ...Retry) (err error) {
	if err = c.SetSendDeadline(time.Now().Add(timeout)); err != nil {
		return err
	}
	defer c.SetSendDeadline(time.Time{})
	err = c.Send(data, retry...)
	return
}

// SendRecv writes data to the connection and blocks reading response.
func (c *Conn) SendRecv(data []byte, length int, retry ...Retry) ([]byte, error) {
	if err := c.Send(data, retry...); err == nil {
		return c.Recv(length, retry...)
	} else {
		return nil, err
	}
}

// SendRecvWithTimeout writes data to the connection and reads response with timeout.
func (c *Conn) SendRecvWithTimeout(data []byte, length int, timeout time.Duration, retry ...Retry) ([]byte, error) {
	if err := c.Send(data, retry...); err == nil {
		return c.RecvWithTimeout(length, timeout, retry...)
	} else {
		return nil, err
	}
}

func (c *Conn) SetDeadline(t time.Time) (err error) {
	if err = c.rawConn.SetDeadline(t); err == nil {
		c.receiveDeadline = t
		c.sendDeadline = t
	}
	if err != nil {
		err = gserror.Wrapf(err, `SetDeadline for connection failed with "%s"`, t)
	}
	return err
}

func (c *Conn) SetReceiveDeadline(t time.Time) (err error) {
	if err = c.rawConn.SetReadDeadline(t); err == nil {
		c.receiveDeadline = t
	}
	if err != nil {
		err = gserror.Wrapf(err, `SetReadDeadline for connection failed with "%s"`, t)
	}
	return err
}

func (c *Conn) SetSendDeadline(t time.Time) (err error) {
	if err = c.rawConn.SetWriteDeadline(t); err == nil {
		c.sendDeadline = t
	}
	if err != nil {
		err = gserror.Wrapf(err, `SetWriteDeadline for connection failed with "%s"`, t)
	}
	return err
}

// SetReceiveBufferWait sets the buffer waiting timeout when reading all data from connection.
// The waiting duration cannot be too long which might delay receiving data from remote address.
func (c *Conn) SetReceiveBufferWait(bufferWaitDuration time.Duration) {
	c.receiveBufferWait = bufferWaitDuration
}

func (c *Conn) Read(b []byte) (n int, err error) {
	return c.reader.Read(b)
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (c *Conn) Write(b []byte) (n int, err error) {
	return c.rawConn.Write(b)
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (c *Conn) Close() error {
	return c.rawConn.Close()
}

// LocalAddr returns the local network address, if known.
func (c *Conn) LocalAddr() net.Addr {
	return c.rawConn.LocalAddr()
}

// RemoteAddr returns the remote network address, if known.
func (c *Conn) RemoteAddr() net.Addr {
	return c.rawConn.RemoteAddr()
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.rawConn.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.rawConn.SetWriteDeadline(t)
}
