package iface

type IRequest interface {
	GetConnection() IConnection
	GetMsg() IMessage
}
