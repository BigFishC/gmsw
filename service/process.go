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
	tag := false
	processes, err := process.Processes()

	if err != nil {
		panic(err)
	}

	for _, p := range processes {
		cmdlines, _ := p.Cmdline()
		if strings.Contains(cmdlines, processname) && !strings.Contains(cmdlines, "gmsf") {
			p.Kill()
			tag = true
			os.Exit(0)
		}
	}
	if tag {
		fmt.Printf("The %s is killed.", processname)
	} else {
		fmt.Printf("The %s is not startted.", processname)
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
	fmt.Printf("The %s service is not running. GO ON!", processname)
	return nil
}
