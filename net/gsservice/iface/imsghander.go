package iface

type IMsgHandler interface {
	Handle(req IRequest) error
}
