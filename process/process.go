package process

import (
	"errors"
	"fmt"
	"log"
	"os"
	"syscall"

	ps "github.com/mitchellh/go-ps"
)

func Process() {
	ps, _ := ps.Processes()
	fmt.Println(ps[0].Executable())

	for pp := range ps {
		fmt.Printf("%d %s\n", ps[pp].Pid(), ps[pp].Executable())
	}
}

// FindProcess( key string ) ( int, string, error )
func FindProcess(key string) (int, string, error) {
	pname := ""
	pid := 0
	err := errors.New("not found")
	ps, _ := ps.Processes()

	for i, _ := range ps {
		if ps[i].Executable() == key {
			pid = ps[i].Pid()
			pname = ps[i].Executable()
			err = nil
			break
		}
	}
	return pid, pname, err
}

func StartDetachedProcess(executable string) {
	err := syscall.Exec(executable, []string{}, os.Environ())
	if err != nil {
		log.Printf("Failed to start process %s: %s", executable, err)
	}
}

func TerminateProcess(pid int, optional_signal ...syscall.Signal) {
	signal := syscall.SIGTERM
	if len(optional_signal) > 0 {
		signal = optional_signal[0]
	}
	if err := syscall.Kill(pid, signal); err != nil {
		log.Fatal("Failed to kill process: ", err)
	}
}
