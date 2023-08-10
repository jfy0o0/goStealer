package gsdns

import (
	"bytes"
	"encoding/binary"
	"github.com/jfy0o0/goStealer/net/gsipv4"
	"golang.org/x/net/dns/dnsmessage"
	"net"
)

func GetDnsUserRequest[T any](src *net.UDPAddr, buf []byte, n int, value T) ([]byte, int, *UserDnsRequest[T]) {
	srcIPString, _, _ := net.SplitHostPort(src.String())

	var (
		srcIP         uint32 = gsipv4.Ip2long(srcIPString)
		springBoardIP uint32
	)

	if n > userProtocolHeaderLen+8 {
		if bytes.Equal(buf[:userProtocolHeaderLen], userProtocolHeader) {
			srcIP = binary.BigEndian.Uint32(buf[userProtocolHeaderLen : userProtocolHeaderLen+4])
			springBoardIP = binary.BigEndian.Uint32(buf[userProtocolHeaderLen+4 : userProtocolHeaderLen+8])
			buf = buf[userProtocolHeaderLen+8:]
			n = n - userProtocolHeaderLen - 8
		}
	}
	msg := &dnsmessage.Message{}
	if err := msg.Unpack(buf[:n]); err != nil {
		return buf, n, nil
	}
	if len(msg.Questions) < 1 {
		return buf, n, nil
	}
	req := NewUserDnsRequest(src, msg, srcIP, springBoardIP, value)
	return buf, n, req
}
