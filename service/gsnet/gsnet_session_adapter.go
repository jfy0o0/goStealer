package gsnet

import "github.com/jfy0o0/goStealer/net/gstcp"

type SessionAdapter[T any] interface {
	OnMsg(conn *gstcp.Conn)
	OnSendMsg(conn *gstcp.Conn, msg interface{})
}
