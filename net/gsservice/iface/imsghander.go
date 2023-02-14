package iface

type IMsgHandler interface {
	HandleOnConnect(conn IConnection)
	HandleCmdChan(req IRequest) error
	HandleDataChan(conn IConnection) error
	HandleOffConnect(conn IConnection)
}
