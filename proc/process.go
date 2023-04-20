package proc

import (
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

//KillProcess func
func KillProcess(processname string) error {
	processes, err := process.Processes()

	if err != nil {
		panic(err)
	}
	for _, p := range processes {
		cmdline, _ := p.Cmdline()
		if strings.Contains(cmdline, processname) {
			p.Kill()
		}
	}
	return nil
}
