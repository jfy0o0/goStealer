package ghssh

import (
	"goHero/net/ghssh/internal"
	"golang.org/x/crypto/ssh"
	"strings"
)

type sshCell struct {
	sshConfig  *ssh.ClientConfig
	remoteAddr string
}

func newByPasswdCellSSH(sshUser, sshPasswd, remoteAddr string) *sshCell {
	if !strings.Contains(remoteAddr, ":") {
		remoteAddr = remoteAddr + ":22"
	}
	sc := &sshCell{
		sshConfig:  internal.GetConfigByPasswdSSH(sshUser, sshPasswd),
		remoteAddr: remoteAddr,
	}
	return sc
}

func newByCertCellSSH(sshUser, sshPrivateKeyPath, remoteAddr string) *sshCell {
	if !strings.Contains(remoteAddr, ":") {
		remoteAddr = remoteAddr + ":22"
	}
	sc := &sshCell{
		sshConfig:  internal.GetConfigByCertSSH(sshUser, sshPrivateKeyPath),
		remoteAddr: remoteAddr,
	}
	return sc
}
