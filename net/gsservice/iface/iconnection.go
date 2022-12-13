package iface

type IConnection interface {
	Start()

	Stop()

	SetFresh(int64)

	GetFresh() int64

	GetConnectionID() string

	SetProperty(string, interface{})

	GetProperty(string) (interface{}, bool)

	DelProperty(string)

	SendMsg(tp byte, data []byte) error
}
