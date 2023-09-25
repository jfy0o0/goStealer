package gsnet

import (
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
	//ctx    context.Context
	//cancel context.CancelFunc

	CommunicationAdapter
}

// , hello *gstcp.GsHello[HelloExtend[T]]
func newServerSession[T any](s *Server[T]) (session *Session[T]) {
	session = &Session[T]{
		Adapter:  s.Config.SessionAdapter,
		Property: gsmap.NewAnyAnyMap[string, interface{}](true),
	}
	session.CommunicationAdapter = GetCommunicationAdapter(s.Config.SessionConf, session)
	return session
}

// , hello *gstcp.GsHello[HelloExtend[T]]
func newClientSession[T any](c *Client[T]) (session *Session[T]) {
	session = &Session[T]{
		Adapter:  c.Config.SessionAdapter,
		Property: gsmap.NewAnyAnyMap[string, interface{}](true),
	}
	session.CommunicationAdapter = GetCommunicationAdapter(c.Config.SessionConf, session)

	return session
}
