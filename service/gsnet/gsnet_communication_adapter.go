package gsnet

import (
	"github.com/jfy0o0/goStealer/net/gstcp"
)

type CommunicationAdapter interface {
	InitSelf(isServer bool, conn *gstcp.Conn) error
	Run()
	Stop()
	Push(interface{})
	CheckHeartBeat() error
}

func GetCommunicationAdapter[T any](config SessionConfig, session *Session[T]) CommunicationAdapter {
	var adapter CommunicationAdapter
	switch config.CommunicationType {
	case CommunicationTypeUserDefined:
		adapter = &CommunicationUserDefined[T]{ParentSession: session,
			Tx: make(chan interface{}, config.TxCap)}
	case CommunicationTypeYamuxSingle:

	case CommunicationTypeYamuxMuti:
		adapter = &CommunicationYamuxMuti[T]{ParentSession: session,
			Tx: make(chan interface{}, config.TxCap)}
	}
	//if err := adapter.InitSelf(isServer); err != nil {
	//	log.Println("can not init self ... ")
	//	return nil
	//}

	return adapter
}
