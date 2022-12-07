package ghssh

import (
	"github.com/pkg/sftp"
	"os"
)

type sshSftp struct {
	*ghssh
	sftpClient *sftp.Client
}

func NewSftpSSH(sshUser, sshPasswd, remoteAddr string) (sf *sshSftp, err error) {
	sf = &sshSftp{}
	sf.ghssh, err = newByPasswdSSH(sshUser, sshPasswd, remoteAddr)
	if err != nil {
		return nil, err
	}
	if err = sf.fillField(); err != nil {
		return nil, err
	}
	return sf, nil
}

func NewSftpByCertSSH(sshUser, sshPrivateKeyPath, remoteAddr string) (sf *sshSftp, err error) {
	sf = &sshSftp{}
	sf.ghssh, err = newByCertSSH(sshUser, sshPrivateKeyPath, remoteAddr)
	if err != nil {
		return nil, err
	}
	if err = sf.fillField(); err != nil {
		return nil, err
	}
	return sf, nil
}

func (sf *sshSftp) fillField() (err error) {
	sf.sftpClient, err = sftp.NewClient(sf.sshClient)
	if err != nil {
		return err
	}
	return nil
}

func (sf *sshSftp) PutFile(sPath, tPath string) error {
	srcFile, err := os.Open(sPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := sf.sftpClient.Create(tPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err = dstFile.ReadFrom(srcFile); err != nil {
		return err
	}
	return nil
}

func (sf *sshSftp) MakeDir(path string) error {
	if _, err := sf.sftpClient.Stat(path); err == nil {
		return nil
	}
	if err := sf.sftpClient.Mkdir(path); err != nil {
		return err
	}
	return nil
}

func (sf *sshSftp) Close() {
	sf.sftpClient.Close()
	sf.ghssh.close()
}
