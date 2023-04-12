package gsdns

import (
	"bytes"
	"encoding/binary"
	"golang.org/x/net/dns/dnsmessage"
	"net"
)

func GetDnsUserRequest[T any](src *net.UDPAddr, buf []byte, n int, value T) ([]byte, int, *UserDnsRequest[T]) {
	var (
		srcIP         uint32
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
