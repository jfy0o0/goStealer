package gstool

import (
	"fmt"
	"os"
	"path"
	"syscall"
)

func GetCurrentProcessName() string {
	return path.Base(os.Args[0])
}

func GetProcessName(pid int) (string, error) {
	pgid, err := syscall.Getpgid(pid)
	if err != nil {
		return "", err
	}
	statFile := fmt.Sprintf("/proc/%d/stat", pgid)
	f, err := os.Open(statFile)
	if err != nil {
		return "", err
	}

	defer f.Close()

	var comm string
	fmt.Fscanf(f, "%d (%s)", &pgid, &comm)
	return comm, nil
}
