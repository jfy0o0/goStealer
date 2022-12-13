package iface

type IMessage interface {
	//Getter
	//GetMsgUuid() uint64
	GetMsgLen() int
	GetMsgData() []byte
	GetMsgTp() byte
	//// Setter
	//SetMsgUuid(uint64)
	//SetMsgId(uint32)
	//SetMsgLen(uint32)
	//SetMsgData([]byte)
}
