package proc

import (
	"github.com/shirou/gopsutil/v3/process"
)

//KillProcess func
func KillProcess(processname string) error {
	processes, err := process.Processes()

	if err != nil {
		panic(err)
	}
	for _, p := range processes {
		n, _ := p.Name()
		if n == processname {
			p.Kill()
		}
	}
	return nil
}
