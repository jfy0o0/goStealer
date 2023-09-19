package gsnet

import (
	"context"
	"github.com/jfy0o0/goStealer/container/gsmap"
	"github.com/jfy0o0/goStealer/net/gstcp"
)

type HelloExtend[T any] struct {
	Key string
	V   T
}
type Session[T any] struct {
	Hello    *gstcp.GsHello[HelloExtend[T]]
	Adapter  SessionAdapter[T]
	Property *gsmap.AnyAnyMap[string, interface{}]
	//Tx           chan interface{}
	// Channel to notify that the connection has exited/stopped
	// (告知该链接已经退出/停止的channel)
	ctx    context.Context
	cancel context.CancelFunc

	CommunicationAdapter
}

// , hello *gstcp.GsHello[HelloExtend[T]]
func newServerSession[T any](s *Server[T]) (session *Session[T]) {
	session = &Session[T]{
		//Hello:    hello,
		Adapter:  s.Config.SessionAdapter,
		Property: gsmap.NewAnyAnyMap[string, interface{}](true),
	}
	//session.CommunicationAdapter = GetCommunicationAdapter(s.Config.SessionConf.CommunicationType, session, true)
	session.CommunicationAdapter = GetCommunicationAdapter(s.Config.SessionConf, session)
	return session
}

// , hello *gstcp.GsHello[HelloExtend[T]]
func newClientSession[T any](c *Client[T]) (session *Session[T]) {
	session = &Session[T]{
		//Hello:    hello,
		Adapter:  c.Config.SessionAdapter,
		Property: gsmap.NewAnyAnyMap[string, interface{}](true),
	}
	session.CommunicationAdapter = GetCommunicationAdapter(c.Config.SessionConf, session)

	return session
}

//func (s *Session[T]) Run() {
//s.communicationAdapter.Run()
//s.ctx, s.cancel = context.WithCancel(context.Background())
//
//go s.runTx()
//
//for {
//	stream, err := s.YamuxSession.Accept()
//	if err != nil {
//		break
//	}
//	go s.handleStream(gstcp.NewConnByNetConn(stream))
//}
//}
//func (s *Session[T]) runTx() {

//for {
//	select {
//	case <-s.ctx.Done():
//		s.YamuxSession.Close()
//		return
//	case msg := <-s.Tx:
//		var err error
//		var c net.Conn
//		for {
//			c, err = s.YamuxSession.Open()
//			if err == nil {
//				break
//			}
//			if s.ctx.Err() != nil {
//				return
//			}
//		}
//
//		s.Adapter.OnSendMsg(gstcp.NewConnByNetConn(c), msg)
//		c.Close()
//	}
//}
//}

//func (s *Session[T]) handleStream(conn *gstcp.Conn) {
//	defer conn.Close()
//	s.Adapter.OnMsg(conn)
//}

//func (s *Session[T]) Stop() {
//s.YamuxSession.Close()
//s.cancel()
//}
