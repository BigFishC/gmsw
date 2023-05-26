package service

import (
	"log"
	"os"
	sh "os/exec"

	"github.com/urfave/cli/v2"
)

type SShell struct {
	PATH string `json:"path"`
	CMD  string `json:"cmd"`
}

func (s *SShell) CheckParam(c *cli.Context) bool {
	if c.String("d") == "" || c.String("c") == "" {
		return false
	} else {
		return true
	}
}

func (s *SShell) StartUp(c *cli.Context) error {

	if c.NArg() == 0 {
		if s.CheckParam(c) {
			s.PATH = c.String("d")

			if err := os.Chdir(s.PATH); err != nil {
				panic(err)
			}
			s.CMD = c.String("c")
			c := sh.Command("bash", "-c", s.CMD)
			if err := c.Start(); err != nil {
				panic(err)
			}
		} else {
			log.Fatal("The parameteris not null!")
		}
	} else {
		log.Fatal("Please use the -h parameter for help")
	}

	return nil
}
