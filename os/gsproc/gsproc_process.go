package gsproc

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const (
	envKeyPPid = "GSPROC_PPID"
)

// Process is the struct for a single process.
type Process struct {
	exec.Cmd
	Manager *Manager
	PPid    int
}

// NewProcess creates and returns a new Process.
func NewProcess(path string, args []string, environment ...[]string) *Process {
	env := os.Environ()
	if len(environment) > 0 {
		for k, v := range environment[0] {
			env[k] = v
		}
	}
	process := &Process{
		Manager: nil,
		PPid:    os.Getpid(),
		Cmd: exec.Cmd{
			Args:       []string{path},
			Path:       path,
			Stdin:      os.Stdin,
			Stdout:     os.Stdout,
			Stderr:     os.Stderr,
			Env:        env,
			ExtraFiles: make([]*os.File, 0),
		},
	}
	process.Dir, _ = os.Getwd()
	if len(args) > 0 {
		// Exclude of current binary path.
		start := 0
		if strings.EqualFold(path, args[0]) {
			start = 1
		}
		process.Args = append(process.Args, args[start:]...)
	}
	return process
}

//// NewProcessCmd creates and returns a process with given command and optional environment variable array.
//func NewProcessCmd(cmd string, environment ...[]string) *Process {
//	return NewProcess(getShell(), append([]string{getShellOption()}, parseCommand(cmd)...), environment...)
//}

// Start starts executing the process in non-blocking way.
// It returns the pid if success, or else it returns an error.
func (p *Process) Start() (int, error) {
	if p.Process != nil {
		return p.Pid(), nil
	}
	p.Env = append(p.Env, fmt.Sprintf("%s=%d", envKeyPPid, p.PPid))
	if err := p.Cmd.Start(); err == nil {
		if p.Manager != nil {
			p.Manager.processes.Set(p.Process.Pid, p)
		}
		return p.Process.Pid, nil
	} else {
		return 0, err
	}
}

// Run executes the process in blocking way.
func (p *Process) Run() error {
	if _, err := p.Start(); err == nil {
		return p.Wait()
	} else {
		return err
	}
}

// Pid retrieves and returns the PID for the process.
func (p *Process) Pid() int {
	if p.Process != nil {
		return p.Process.Pid
	}
	return 0
}

//// Send sends custom data to the process.
//func (p *Process) Send(data []byte) error {
//	if p.Process != nil {
//		return Send(p.Process.Pid, data)
//	}
//	return gserror.NewCode(gscode.CodeInvalidParameter, "invalid process")
//}

// Release releases any resources associated with the Process p,
// rendering it unusable in the future.
// Release only needs to be called if Wait is not.
func (p *Process) Release() error {
	return p.Process.Release()
}

// Kill causes the Process to exit immediately.
func (p *Process) Kill() error {
	if err := p.Process.Kill(); err == nil {
		if p.Manager != nil {
			p.Manager.processes.Remove(p.Pid())
		}
		if runtime.GOOS != "windows" {
			if err = p.Process.Release(); err != nil {
			}
		}
		_, err = p.Process.Wait()
		//return err
		return nil
	} else {
		return err
	}
}

// Signal sends a signal to the Process.
// Sending Interrupt on Windows is not implemented.
func (p *Process) Signal(sig os.Signal) error {
	return p.Process.Signal(sig)
}
