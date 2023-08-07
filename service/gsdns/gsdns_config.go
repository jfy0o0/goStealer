package gsdns

import "github.com/jfy0o0/goStealer/container/gstype"

var innerID = gstype.NewInt(0)

type TxServerConfig[T any] struct {
	Id           int
	ToDnsIP      string
	AppendHeader bool
	LocalIP      uint32
	ScanInterval int
	OnNoResponse func(*UserDnsRequest[T])
}

func GetDefaultTxServerConfig[T any]() *TxServerConfig[T] {
	return &TxServerConfig[T]{
		Id:           innerID.Add(1),
		ToDnsIP:      "8.8.8.8:53",
		AppendHeader: false,
		LocalIP:      0,
		ScanInterval: 5,
		OnNoResponse: nil,
	}
}
