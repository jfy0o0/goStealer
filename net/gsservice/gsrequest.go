package gsservice

import "github.com/jfy0o0/goStealer/net/gsservice/iface"

type Request struct {
	conn iface.IConnection
	msg  iface.IMessage
}

func NewRequest(conn iface.IConnection, msg iface.IMessage) iface.IRequest {
	return &Request{
		conn: conn,
		msg:  msg,
	}
}

func (r *Request) GetConnection() iface.IConnection {
	return r.conn
}

func (r *Request) GetMsg() iface.IMessage {
	return r.msg
}
