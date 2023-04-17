package util

import (
	"fmt"
	"log"
	"strings"

	"github.com/BigFishC/gmsw/config"
	"github.com/BigFishC/gmsw/secret"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh"
)

type Cli struct {
	USER      string      `json:"user"`
	PWD       string      `json:"pwd"`
	IP        string      `json:"ip"`
	PORT      string      `json:"port"`
	SSHCLIENT *ssh.Client `json:"sshclient"`
}

func (c *Cli) getConfig() *ssh.ClientConfig {
	config := &ssh.ClientConfig{
		User: c.USER,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.PWD),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return config
}

//Connect
func (c *Cli) Connect() error {
	client, err := ssh.Dial("tcp", c.IP+":"+c.PORT, c.getConfig())
	if err != nil {
		panic(err)
	}
	c.SSHCLIENT = client
	// defer client.Close()
	return nil
}

//ChangeEnv
func (c *Cli) ChangeEnv(envparam string, pwdparam string, cli *cli.Context) error {
	decrypt, err := secret.DecryptByAes(pwdparam, secret.PwdKey)
	if err != nil {
		panic(err)
	}
	analysiInfo := cli.String(envparam)
	fmt.Println(analysiInfo)
	analysiStringSplite := strings.Split(analysiInfo, "@")
	c.USER = analysiStringSplite[0]
	c.IP = analysiStringSplite[1]
	c.PWD = string(decrypt)

	return nil
}

//Run
func (c *Cli) Run(cli *cli.Context) error {
	var config config.ConfigStruct
	config.LoadConfig()
	analysiCmd := cli.Args().Get(0)

	if cli.String("P") == "" {
		env := cli.FlagNames()[0]
		switch env {
		case "t":
			c.ChangeEnv("t", config.Tpwd, cli)
		case "p":
			c.ChangeEnv("p", config.Ppwd, cli)
		}
		c.PORT = "22"
	} else {
		env := cli.FlagNames()[1]
		switch env {
		case "t":
			c.ChangeEnv("t", config.Tpwd, cli)
		case "p":
			c.ChangeEnv("p", config.Ppwd, cli)
		}
		c.PORT = cli.String("P")
	}

	if cli.Args().Len() > 0 {
		if err := c.Connect(); err != nil {
			log.Fatal(err)
		}
		session, err := c.SSHCLIENT.NewSession()
		if err != nil {
			log.Fatal(err)
		}
		defer session.Close()
		if err := session.Run(analysiCmd); err != nil {
			log.Fatalf("Failed to run %s", analysiCmd)
		}
		log.Fatalf("Successed to run %s", analysiCmd)
	} else {
		log.Fatal("Param is too match!")
	}

	return nil
}
