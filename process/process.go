package process

import (
	"errors"
	"fmt"
	"os"
	"strings"

	ps "github.com/mitchellh/go-ps"
)

type Process struct {
	Pid  int
	Name string
}

// FindProcess( key string ) ( int, string, error )
func FindProcess(key string) ([]Process, error) {
	err := errors.New("not found")
	ps, _ := ps.Processes()
	var processes []Process

	for i := range ps {
		if strings.Contains(ps[i].Executable(), key) {
			processes = append(processes, Process{ps[i].Pid(), ps[i].Executable()})
			err = nil
		}
	}
	return processes, err
}

func KillProcesses(procs []Process) error {
	var killError error
	for _, proc := range procs {
		if p, _ := GetProcess(proc.Pid); p.Pid != 0 {
			if err := p.Signal(os.Kill); err != nil {
				fmt.Printf("Unable to kill process [%s]: %s", proc.Name, err)
				killError = err
				return err
			}
		}
	}
	return killError
}

func GetProcess(pid int) (*os.Process, error) {
	return os.FindProcess(pid)
}
