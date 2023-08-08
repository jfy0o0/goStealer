package gsdns

import (
	"github.com/jfy0o0/goStealer/container/gstype"
	"golang.org/x/net/dns/dnsmessage"
	"net"
	"time"
)

type UserDnsRequest[T any] struct {
	addr          *net.UDPAddr
	msg           *dnsmessage.Message
	endTime       int64
	oldDnsID      uint16
	count         *gstype.Int
	srcIP         uint32
	springBoardIP uint32
	value         T
}

func NewUserDnsRequest[T any](addr *net.UDPAddr, msg *dnsmessage.Message, srcIP, springBoardIP uint32, value T) *UserDnsRequest[T] {
	return &UserDnsRequest[T]{
		addr:          addr,
		msg:           msg,
		endTime:       time.Now().Add(time.Second * 5).Unix(),
		oldDnsID:      msg.ID,
		value:         value,
		count:         gstype.NewInt(),
		srcIP:         srcIP,
		springBoardIP: springBoardIP,
	}
}

func (r *UserDnsRequest[T]) GetDnsMessage() *dnsmessage.Message {
	return r.msg
}
func (r *UserDnsRequest[T]) GetRequest() *net.UDPAddr {
	return r.addr
}
func (r *UserDnsRequest[T]) GetSrcIP() uint32 {
	return r.srcIP
}
func (r *UserDnsRequest[T]) GetSpringBoardIP() uint32 {
	return r.springBoardIP
}
func (r *UserDnsRequest[T]) GetValue() T {
	return r.value
}

func (r *UserDnsRequest[T]) SetValue(value T) {
	r.value = value
}
