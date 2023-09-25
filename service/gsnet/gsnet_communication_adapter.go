package gsnet

import (
	"github.com/jfy0o0/goStealer/net/gstcp"
	"log"
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
		adapter = NewCommunicationUserDefinedFromConfig[T](config, session)
	case CommunicationTypeYamuxSingle:

	case CommunicationTypeYamuxMuti:
		log.Fatalln("not support ")
		adapter = NewCommunicationYamuxMutiFromConfig[T](config, session)
	}
	//if err := adapter.InitSelf(isServer); err != nil {
	//	log.Println("can not init self ... ")
	//	return nil
	//}

	return adapter
}
