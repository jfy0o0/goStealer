package internal

import (
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"time"
)

func GetConfigByPasswdSSH(sshUser, sshPasswd string) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User:            sshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshPasswd),
		},
	}
}

func GetConfigByCertSSH(sshUser, sshPrivateKeyPath string) *ssh.ClientConfig {
	config := &ssh.ClientConfig{
		User:            sshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}
	key, err := os.ReadFile(sshPrivateKeyPath)
	if err != nil {
		log.Fatalln(err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalln(err)
	}
	config.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	return config
}
