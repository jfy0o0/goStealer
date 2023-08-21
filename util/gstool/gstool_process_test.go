package gstool

import (
	"os"
	"testing"
)

func TestGetCurrentProcessName(t *testing.T) {
	name := GetCurrentProcessName()
	t.Logf("name : %v", name)

}
func TestGetProcessName(t *testing.T) {
	pid := os.Getpid()
	t.Logf("pid : %v", pid)
	name, err := GetProcessName(pid)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(name)
}
