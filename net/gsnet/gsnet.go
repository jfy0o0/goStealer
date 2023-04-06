package gsnet

import (
	"errors"
	"golang.org/x/sys/unix"
	"net"
	"syscall"
)

func GetRawDstIpPort(conn net.Conn) (string, int, error) {

	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return "", 0, errors.New("assert [*net.TCPConn] failed ")
	}
	tcpConnFile, err := tcpConn.File()
	if err != nil {
		return "", 0, err
	}
	defer tcpConnFile.Close()
	addr, err := syscall.GetsockoptIPv6Mreq(int(tcpConnFile.Fd()), syscall.SOL_IP, unix.SO_ORIGINAL_DST)
	if err != nil {
		return "", 0, err
	}
	port := int(addr.Multiaddr[2])
	port = port<<8 + int(addr.Multiaddr[3])
	ip := net.IPv4(addr.Multiaddr[4], addr.Multiaddr[5], addr.Multiaddr[6], addr.Multiaddr[7]).String()
	return ip, port, nil
}
