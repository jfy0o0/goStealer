package gsdns

import (
	"golang.org/x/net/dns/dnsmessage"
	"net"
)

type UserDnsResponse[T any] struct {
	msg   *dnsmessage.Message
	addr  *net.UDPAddr
	value T
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
