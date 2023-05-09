package util

import (
	"fmt"
	"log"
	"net"
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
	TRANSFER   string       `json:"transfer"`
	ENV        string       `json:"enviroment"`
	EFLAG      int          `json:"eflag"`
	SOURCEFILE string       `json:"sourcefile"`
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

// Check host parameter
func CheckHost(host string) bool {
	infos := strings.Split(host, "@")
	ip := infos[1]
	if net.ParseIP(ip) != nil {
		return true
	}
	return false
}

// ChangeEnv
func (c *Cli) ChangeEnv(envparam string, pwdparam string, cli *cli.Context) error {
	decrypt, err := secret.DecryptByAes(pwdparam, secret.PwdKey)
	if err != nil {
		panic(err)
	}
	analysiInfo := cli.String(envparam)
	if analysiInfo == "" || CheckHost(analysiInfo) == false {
		log.Fatal("The parameter is error !")
	} else {
		analysiStringSplite := strings.Split(analysiInfo, "@")
		c.USER = analysiStringSplite[0]
		c.IP = analysiStringSplite[1]
		c.PWD = string(decrypt)
	}
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
		if err != nil {
			panic(err)
		}
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
		for _, v := range cli.FlagNames() {
			switch v {
			case "P":
				if cli.String("P") == "" {
					c.PORT = "22"
				} else {
					c.PORT = cli.String("P")
				}
			case "T":
				c.TRANSFER = "T"
				if cli.String("T") == "" {
					log.Fatal("Source file could not be null !")
				} else {
					c.SOURCEFILE = cli.String("T")
				}
			case "t":
				c.ENV = "t"
				c.EFLAG = c.EFLAG + 1
			case "p":
				c.ENV = "p"
				c.EFLAG = c.EFLAG + 1
			}
		}
		c.ChangeEnv(c.ENV, config.Tpwd, cli)
		analysiCmd := cli.Args().Get(0)
		if c.PORT == "" {
			c.PORT = "22"
		}
		if c.EFLAG == 1 {
			if c.TRANSFER == "" {
				c.Run(analysiCmd, cli)
			} else {
				c.UploadFile(c.SOURCEFILE, analysiCmd, cli)
			}
		} else {
			log.Fatal("Enviroment is error !")
		}
	} else {
		log.Fatal("Please use the -h parameter for help")
	}
	return nil
}
