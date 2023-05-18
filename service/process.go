package service

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

// KillProcess func
func KillProcess(processname string) error {
	processes, err := process.Processes()

	if err != nil {
		panic(err)
	}
	for _, p := range processes {
		cmdline, _ := p.Cmdline()
		if strings.Contains(cmdline, processname) && !strings.Contains(cmdline, "gmsf") {
			p.Kill()
			os.Exit(0)
		} else {
			log.Fatalf("%s is not exist", processname)
		}
	}
	return nil
}

func ProcessStatus(processname string) error {
	processes, err := process.Processes()

	if err != nil {
		panic(err)
	}
	for _, p := range processes {
		cmdline, _ := p.Cmdline()
		if strings.Contains(cmdline, processname) && !strings.Contains(cmdline, "gmsf") {
			log.Fatalf("%s is exist! Please check it!", processname)
		}
	}
	fmt.Printf("%s is not exist! Go on!\n", processname)
	return nil
}
