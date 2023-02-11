package iface

type IMsgHandler interface {
	HandleCmdChan(req IRequest) error
	HandleDataChan(conn IConnection)
}
