package gsdns

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"github.com/jfy0o0/goStealer/net/gsipv4"
	"golang.org/x/net/dns/dnsmessage"
	"net"
	"os"
	"strings"
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

func GetLocalDnsIPs() (ips []string) {
	f, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}

		if !strings.HasPrefix(string(line), "nameserver") {
			continue
		}

		_, after, ok := strings.Cut(string(line), " ")
		if !ok {
			continue
		}
		ips = append(ips, after)
	}
	return
}
