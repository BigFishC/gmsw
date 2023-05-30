package service

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
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
	if len(infos) != 2 {
		return false
	} else if net.ParseIP(infos[1]) != nil {
		return true
	} else {
		return false
	}

}

// ChangeEnv
func (c *Cli) ChangeEnv(envparam string, pwdparam string, cli *cli.Context) error {
	decrypt, err := secret.DecryptByAes(pwdparam, secret.PwdKey)
	if err != nil {
		panic(err)
	}
	analysiInfo := cli.String(envparam)
	if analysiInfo == "" || !CheckHost(analysiInfo) {
		log.Fatal("The parameter is error !")
	} else {
		analysiStringSplite := strings.Split(analysiInfo, "@")
		c.USER = analysiStringSplite[0]
		c.IP = analysiStringSplite[1]
		c.PWD = string(decrypt)
	}
	return nil
}

func (c *Cli) UploadFile(localFile string, remoteFile string, cli *cli.Context) error {

	if cli.Args().Len() > 0 {
		if err := c.Connect(); err != nil {
			log.Fatal(err)
		}
		file, err := os.Open(localFile)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		fileState, err := file.Stat()
		if fileState.IsDir() {
			if err := c.SFTPCLIENT.MkdirAll(remoteFile); err != nil {
				log.Fatal(err)
			}
			localFiles, _ := ioutil.ReadDir(localFile)
			for _, localf := range localFiles {
				nextRemoteFilePath := path.Join(remoteFile, localf.Name())
				nextLocalFilePath := path.Join(localFile, localf.Name())
				c.UploadFile(nextLocalFilePath, nextRemoteFilePath, cli)
			}
		} else {
			ftpFile, err := c.SFTPCLIENT.Create(remoteFile)
			if err != nil {
				log.Fatal(err)
			}
			defer ftpFile.Close()
			if err != nil {
				panic(err)
			}
			if _, err := ftpFile.ReadFromWithConcurrency(file, 5); err != nil {
				panic(err)
			}
			fmt.Printf("Successed to transfer %s;\n", remoteFile)
		}

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
				c.PWD = config.Tpwd
			case "p":
				c.ENV = "p"
				c.EFLAG = c.EFLAG + 1
				c.PWD = config.Ppwd
			}
		}
		c.ChangeEnv(c.ENV, c.PWD, cli)
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
