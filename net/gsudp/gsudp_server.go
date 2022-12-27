package gsudp

import (
	"github.com/jfy0o0/goStealer/errors/gscode"
	"github.com/jfy0o0/goStealer/errors/gserror"
	"net"
)

// Server is the UDP server.
type Server struct {
	conn    *Conn       // UDP server connection object.
	address string      // UDP server listening address.
	handler func(*Conn) // Handler for UDP connection.
}

// NewServer creates and returns a UDP server.
// The optional parameter `name` is used to specify its name, which can be used for
// GetServer function to retrieve its instance.
func NewServer(address string, handler func(*Conn)) *Server {
	s := &Server{
		address: address,
		handler: handler,
	}
	return s
}

// SetAddress sets the server address for UDP server.
func (s *Server) SetAddress(address string) {
	s.address = address
}

// SetHandler sets the connection handler for UDP server.
func (s *Server) SetHandler(handler func(*Conn)) {
	s.handler = handler
}

// Close closes the connection.
// It will make server shutdowns immediately.
func (s *Server) Close() (err error) {
	err = s.conn.Close()
	if err != nil {
		err = gserror.Wrap(err, "connection failed")
	}
	return
}

// Run starts listening UDP connection.
func (s *Server) Run() error {
	if s.handler == nil {
		err := gserror.NewCode(gscode.CodeMissingConfiguration, "start running failed: socket handler not defined")
		return err
	}
	addr, err := net.ResolveUDPAddr("udp", s.address)
	if err != nil {
		err = gserror.Wrapf(err, `net.ResolveUDPAddr failed for address "%s"`, s.address)
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		err = gserror.Wrapf(err, `net.ListenUDP failed for address "%s"`, s.address)
		return err
	}
	s.conn = NewConnByNetConn(conn)
	s.handler(s.conn)
	return nil
}
