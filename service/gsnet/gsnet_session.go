package gsnet

import (
	"context"
	"github.com/hashicorp/yamux"
	"github.com/jfy0o0/goStealer/container/gsmap"
	"github.com/jfy0o0/goStealer/net/gstcp"
)

type HelloExtend[T any] struct {
	Key string
	V   T
}
type Session[T any] struct {
	YamuxSession *yamux.Session
	Hello        *gstcp.GsHello[HelloExtend[T]]
	Adapter      SessionAdapter[T]
	Property     *gsmap.AnyAnyMap[string, interface{}]
	Tx           chan interface{}
	// Channel to notify that the connection has exited/stopped
	// (告知该链接已经退出/停止的channel)
	ctx    context.Context
	cancel context.CancelFunc
}

func newServerSession[T any](s *Server[T], conn *gstcp.Conn, hello *gstcp.GsHello[HelloExtend[T]]) (session *Session[T], err error) {
	session = &Session[T]{
		Hello:    hello,
		Adapter:  s.Config.SessionAdapter,
		Property: gsmap.NewAnyAnyMap[string, interface{}](true),
		Tx:       make(chan interface{}, 1024),
	}
	session.YamuxSession, err = yamux.Server(conn, nil)
	return session, err
}

func newClientSession[T any](c *Client[T], conn *gstcp.Conn, hello *gstcp.GsHello[HelloExtend[T]]) (session *Session[T], err error) {
	session = &Session[T]{
		Hello:    hello,
		Adapter:  c.Config.SessionAdapter,
		Property: gsmap.NewAnyAnyMap[string, interface{}](true),
		Tx:       make(chan interface{}, 1024),
	}
	session.YamuxSession, err = yamux.Client(conn, nil)
	return session, err
}

func (s *Session[T]) Run() {
	s.ctx, s.cancel = context.WithCancel(context.Background())

	go s.runTx()

	for {
		stream, err := s.YamuxSession.Accept()
		if err != nil {
			break
		}
		go s.handleStream(gstcp.NewConnByNetConn(stream))
	}
}
func (s *Session[T]) runTx() {
	for {
		select {
		case <-s.ctx.Done():
			s.YamuxSession.Close()
			return
		case msg := <-s.Tx:
			c, err := s.YamuxSession.Open()
			if err != nil {
				continue
			}
			s.Adapter.OnSendMsg(gstcp.NewConnByNetConn(c), msg)
		}
	}
}

func (s *Session[T]) handleStream(conn *gstcp.Conn) {
	defer conn.Close()
	s.Adapter.OnMsg(conn)
}

func (s *Session[T]) Stop() {
	s.cancel()
}
