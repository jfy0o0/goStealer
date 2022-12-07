package gsssh

import (
	"golang.org/x/crypto/ssh"
)

type gsssh struct {
	sshClient *ssh.Client
	*sshCell
}

func newByPasswdSSH(sshUser, sshPasswd, remoteAddr string) (gs *gsssh, err error) {
	gs = &gsssh{
		sshCell: newByPasswdCellSSH(sshUser, sshPasswd, remoteAddr),
	}
	gs.sshClient, err = ssh.Dial("tcp", gs.remoteAddr, gs.sshConfig)
	if err != nil {
		return nil, err
	}
	return gs, nil
}

func newByCertSSH(sshUser, sshPrivateKeyPath, remoteAddr string) (gs *gsssh, err error) {
	gs = &gsssh{
		sshCell: newByCertCellSSH(sshUser, sshPrivateKeyPath, remoteAddr),
	}
	gs.sshClient, err = ssh.Dial("tcp", gs.remoteAddr, gs.sshConfig)
	if err != nil {
		return nil, err
	}
	return gs, nil
}

func (gs *gsssh) RunCmd(cmd string, sync bool) error {
	session, err := gs.sshClient.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	if sync {
		return session.Run(cmd)
	}
	return session.Start(cmd)
}

func (gs *gsssh) close() {
	gs.sshClient.Close()
}
