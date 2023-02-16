package iface

import "github.com/jfy0o0/goStealer/net/gstcp"

type IConnection interface {
	Start()

	Stop()

	SetFresh(int64)

	GetFresh() int64

	GetConnectionID() string

	SetProperty(string, interface{})

	GetProperty(string) (interface{}, bool)

	DelProperty(string)

	SendLvMsg(data []byte)

	IsCmdChan() bool

	GetRawConn() *gstcp.Conn
}
