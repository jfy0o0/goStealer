package gshost

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"sync"
)

const (
	initGsHostBufSize = 1024 // allocate 1 KB up front to try to avoid resizing
)

type SharedConn struct {
	sync.Mutex
	net.Conn
	gsHostBuf *bytes.Buffer // all of the initial data that has to be read in order to vhost a connection is saved here
}

func NewShared(conn net.Conn, f ...func(tempReader *bufio.Reader)) (*SharedConn, io.Reader) {
	c := &SharedConn{
		Conn:      conn,
		gsHostBuf: bytes.NewBuffer(make([]byte, 0, initGsHostBufSize)),
	}
	reader := bufio.NewReader(conn)

	if len(f) > 0 {
		f[0](reader)
	}
	return c, io.TeeReader(reader, c.gsHostBuf)
}

func (c *SharedConn) Read(p []byte) (n int, err error) {
	c.Lock()
	if c.gsHostBuf == nil {
		c.Unlock()
		return c.Conn.Read(p)
	}
	n, err = c.gsHostBuf.Read(p)

	// end of the request buffer
	if err == io.EOF {
		// let the request buffer get garbage collected
		// and make sure we don't read from it again
		c.gsHostBuf = nil

		// continue reading from the connection
		var n2 int
		n2, err = c.Conn.Read(p[n:])

		// update total read
		n += n2
	}
	c.Unlock()
	return
}
