package ghssh

import (
	"fmt"
	"goHero/errors/gherror"
	"goHero/net/ghssh/internal"
	"golang.org/x/crypto/ssh"
	"net"
)

type sshTunnel struct {
	sshCell   *sshCell
	localPort string
	destAddr  string
}

func NewTunnelSSH(sshUser, sshPasswd, localPort, remoteAddr, destAddr string) *sshTunnel {
	st := &sshTunnel{
		sshCell:   newByPasswdCellSSH(sshUser, sshPasswd, remoteAddr),
		localPort: localPort,
		destAddr:  destAddr,
	}
	return st
}

func NewTunnelByCertSSH(sshUser, sshPrivateKeyPath, localPort, remoteAddrSSH, destAddr string) *sshTunnel {
	st := &sshTunnel{
		sshCell:   newByCertCellSSH(sshUser, sshPrivateKeyPath, remoteAddrSSH),
		localPort: localPort,
		destAddr:  destAddr,
	}
	return st
}

func (st *sshTunnel) Run() error {
	listen, err := net.Listen("tcp", ":"+st.localPort)
	if err != nil {
		return gherror.Wrapf(err, "listen port [%v] failed", st.localPort)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed error :", err)
			continue
		}
		go func() {
			if err := st.forward(conn); err != nil {
				fmt.Println("forward error: ", err, conn.RemoteAddr())
			} else {
				fmt.Println("forward close :", conn.RemoteAddr())
			}
		}()
	}
}

func (st *sshTunnel) forward(localConn net.Conn) error {
	fmt.Println("new conn :", localConn.RemoteAddr())
	defer localConn.Close()

	sshClientConn, err := ssh.Dial("tcp", st.sshCell.remoteAddr, st.sshCell.sshConfig)
	if err != nil {
		return err
	}
	defer sshClientConn.Close()

	sshConn, err := sshClientConn.Dial("tcp", st.destAddr)
	if err != nil {
		return err
	}
	defer sshConn.Close()

	internal.Relay(localConn, sshConn)

	return nil
}
