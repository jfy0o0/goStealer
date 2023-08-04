package gstcp

import (
	"crypto/tls"
	"github.com/jfy0o0/goStealer/container/gsmap"
	"github.com/jfy0o0/goStealer/errors/gscode"
	"github.com/jfy0o0/goStealer/errors/gserror"
	"log"
	"net"
	"sync"
)

const (
	// defaultServer is the default TCP server name.
	defaultServer = "default"
)

// Server is a TCP server.
type Server struct {
	mu        sync.Mutex   // Used for Server.listen concurrent safety.
	listen    net.Listener // Listener.
	address   string       // Server listening address.
	handler   func(*Conn)  // Connection handler.
	tlsConfig *tls.Config  // TLS configuration.
}

var serverMapping = gsmap.NewAnyAnyMap[string, *Server](true)

func GetServer(name ...string) *Server {
	serverName := defaultServer
	if len(name) > 0 && name[0] != "" {
		serverName = name[0]
	}
	return serverMapping.GetOrSetFunc(serverName, func() *Server {
		return NewServer("", nil)
	})
}

// NewServer creates and returns a new normal TCP server.
// The parameter <name> is optional, which is used to specify the instance name of the server.
func NewServer(address string, handler func(*Conn), name ...string) *Server {
	s := &Server{
		address: address,
		handler: handler,
	}
	if len(name) > 0 && name[0] != "" {
		serverMapping.Set(name[0], s)
	}
	return s
}

// NewServerTLS creates and returns a new TCP server with TLS support.
// The parameter <name> is optional, which is used to specify the instance name of the server.
func NewServerTLS(address string, tlsConfig *tls.Config, handler func(*Conn), name ...string) *Server {
	s := NewServer(address, handler, name...)
	s.SetTLSConfig(tlsConfig)
	return s
}

// NewServerKeyCrt creates and returns a new TCP server with TLS support.
// The parameter <name> is optional, which is used to specify the instance name of the server.
func NewServerKeyCrt(address, crtFile, keyFile string, handler func(*Conn), name ...string) *Server {
	s := NewServer(address, handler, name...)
	if err := s.SetTLSKeyCrt(crtFile, keyFile); err != nil {
		log.Println(err)
	}
	return s
}

// SetAddress sets the listening address for server.
func (s *Server) SetAddress(address string) {
	s.address = address
}

// SetHandler sets the connection handler for server.
func (s *Server) SetHandler(handler func(*Conn)) {
	s.handler = handler
}

// SetTLSKeyCrt sets the certificate and key file for TLS configuration of server.
func (s *Server) SetTLSKeyCrt(crtFile, keyFile string) error {
	tlsConfig, err := LoadKeyCrt(crtFile, keyFile)
	if err != nil {
		return err
	}
	s.tlsConfig = tlsConfig
	return nil
}

// SetTLSConfig sets the TLS configuration of server.
func (s *Server) SetTLSConfig(tlsConfig *tls.Config) {
	s.tlsConfig = tlsConfig
}

// Close closes the listener and shutdowns the server.
func (s *Server) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.listen == nil {
		return nil
	}
	return s.listen.Close()
}

// Run starts running the TCP Server.
func (s *Server) Run() (err error) {
	if s.handler == nil {
		err = gserror.NewCode(gscode.CodeMissingConfiguration, "start running failed: socket handler not defined")
		log.Println(err)
		return
	}
	if s.tlsConfig != nil {
		// TLS Server
		s.mu.Lock()
		s.listen, err = tls.Listen("tcp", s.address, s.tlsConfig)
		s.mu.Unlock()
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		// Normal Server
		addr, err := net.ResolveTCPAddr("tcp", s.address)
		if err != nil {
			log.Println(err)
			return err
		}
		s.mu.Lock()
		s.listen, err = net.ListenTCP("tcp", addr)
		s.mu.Unlock()
		if err != nil {
			log.Println(err)
			return err
		}
	}
	// Listening loop.
	for {
		if conn, err := s.listen.Accept(); err != nil {
			return err
		} else if conn != nil {
			go s.handler(NewConnByNetConn(conn))
		}
	}
}
