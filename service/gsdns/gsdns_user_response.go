package gsdns

import (
	"golang.org/x/net/dns/dnsmessage"
	"net"
)

type UserDnsResponse[T any] struct {
	msg           *dnsmessage.Message
	addr          *net.UDPAddr
	resResponseIP uint32
	value         T
}

func NewUserDnsResponse[T any](msg *dnsmessage.Message, addr *net.UDPAddr, resResponseIP uint32, value T) *UserDnsResponse[T] {
	return &UserDnsResponse[T]{
		msg:           msg,
		addr:          addr,
		resResponseIP: resResponseIP,
		value:         value,
	}
}
func (r *UserDnsResponse[T]) GetResponseIP() uint32 {
	return r.resResponseIP
}
func (r *UserDnsResponse[T]) GetDnsMessage() *dnsmessage.Message {
	return r.msg
}

func (r *UserDnsResponse[T]) GetDestAddr() *net.UDPAddr {
	return r.addr
}

func (r *UserDnsResponse[T]) GetValue() T {
	return r.value
}
