package util

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BigFishC/gmsw/config"
	"github.com/BigFishC/gmsw/secret"
	"github.com/pkg/sftp"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh"
)

type Cli struct {
	USER       string       `json:"user"`
	PWD        string       `json:"pwd"`
	IP         string       `json:"ip"`
	PORT       string       `json:"port"`
	SSHCLIENT  *ssh.Client  `json:"sshclient"`
	SFTPCLIENT *sftp.Client `json:"sftpclient"`
}

const size = 1 << 15

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

// Connect
func (c *Cli) Connect() error {
	client, err := ssh.Dial("tcp", c.IP+":"+c.PORT, c.getConfig())
	if err != nil {
		panic(err)
	}
	sftp, err := sftp.NewClient(client, sftp.MaxPacket(size))
	if err != nil {
		panic(err)
	}
	c.SSHCLIENT = client
	c.SFTPCLIENT = sftp
	return nil
}

// ChangeEnv
func (c *Cli) ChangeEnv(envparam string, pwdparam string, cli *cli.Context) error {
	decrypt, err := secret.DecryptByAes(pwdparam, secret.PwdKey)
	if err != nil {
		panic(err)
	}
	analysiInfo := cli.String(envparam)
	analysiStringSplite := strings.Split(analysiInfo, "@")
	c.USER = analysiStringSplite[0]
	c.IP = analysiStringSplite[1]
	c.PWD = string(decrypt)

	return nil
}

func (c *Cli) UploadFile(localfile string, remotefile string, cli *cli.Context) error {
	file, err := os.Open(localfile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if cli.Args().Len() > 0 {
		if err := c.Connect(); err != nil {
			log.Fatal(err)
		}

		ftpFile, err := c.SFTPCLIENT.Create(remotefile)

		if err != nil {
			log.Fatal(err)
		}
		defer ftpFile.Close()

		// fileByte, err := ioutil.ReadAll(file)
		// fileByte, err := io.ReadFull(file, make([]byte, 1e9))
		if err != nil {
			panic(err)
		}

		// if _, err := ftpFile.Write(fileByte); err != nil {
		if _, err := ftpFile.ReadFromWithConcurrency(file, 10); err != nil {
			panic(err)
		}
		fmt.Printf("Successed to transfer %s", remotefile)
	} else {
		log.Fatal("Param is too match!")
	}
	return nil
}

func (c *Cli) Run(cmd string, cli *cli.Context) error {
	if cli.Args().Len() > 0 {
		if err := c.Connect(); err != nil {
			log.Fatal(err)
		}
		session, err := c.SSHCLIENT.NewSession()
		if err != nil {
			log.Fatal(err)
		}
		defer session.Close()
		if err := session.Run(cmd); err != nil {
			log.Fatalf("Failed to run %s", cmd)
		}
		fmt.Printf("Successed to run %s\n", cmd)

	} else {
		log.Fatal("Param is too match!")
	}
	return nil
}

// Run server
func (c *Cli) Server(cli *cli.Context) error {

	var config config.ConfigStruct
	config.LoadConfig()
	if cli.NArg() > 0 {
		analysiCmd := cli.Args().Get(0)
		if cli.String("P") == "" {
			c.PORT = "22"
			if cli.FlagNames()[0] == "T" {
				env := cli.FlagNames()[1]
				switch env {
				case "t":
					c.ChangeEnv("t", config.Tpwd, cli)
				case "p":
					c.ChangeEnv("p", config.Ppwd, cli)
				}
				localFile := cli.String("T")
				c.UploadFile(localFile, analysiCmd, cli)
			} else {
				env := cli.FlagNames()[0]
				switch env {
				case "t":
					c.ChangeEnv("t", config.Tpwd, cli)
				case "p":
					c.ChangeEnv("p", config.Ppwd, cli)
				}
				c.Run(analysiCmd, cli)
			}

		} else {
			c.PORT = cli.String("P")
			if cli.FlagNames()[1] == "T" {
				env := cli.FlagNames()[2]
				switch env {
				case "t":
					c.ChangeEnv("t", config.Tpwd, cli)
				case "p":
					c.ChangeEnv("p", config.Ppwd, cli)
				}
				localFile := cli.String("T")
				c.UploadFile(localFile, analysiCmd, cli)
			} else {
				env := cli.FlagNames()[1]
				switch env {
				case "t":
					c.ChangeEnv("t", config.Tpwd, cli)
				case "p":
					c.ChangeEnv("p", config.Ppwd, cli)
				}
				c.Run(analysiCmd, cli)
			}

		}
	} else {
		log.Fatal("Please use the -h parameter for help")
	}

	return nil
}
