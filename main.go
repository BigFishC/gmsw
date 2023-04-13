package main

import (
	"os"

	"github.com/BigFishC/gmsw/config"
	"github.com/BigFishC/gmsw/util"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{}
	app.Name = "gmsw"
	app.Usage = "golang middle software"
	app.Version = "1.0.0"
	app.Commands = []*cli.Command{
		Encrypt(),
		RunCmd(),
	}
	// var config config.ConfigStruct
	// config.LoadConfig()

	// decrypt, _ := secret.DecryptByAes(config.Testpwd, secret.PwdKey)

	// fmt.Printf("解密后：%s\n", decrypt)
	app.Run(os.Args)
}

func Encrypt() *cli.Command {
	return &cli.Command{
		Name:   "encrypt",
		Usage:  "encrypt --tpwd=string | encrypt --ppwd=string",
		Action: (&config.ConfigStruct{}).UpdateConfig,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "tpwd",
				Usage: "--tpwd",
			},
			&cli.StringFlag{
				Name:  "ppwd",
				Usage: "--ppwd",
			},
		},
	}
}

func RunCmd() *cli.Command {
	return &cli.Command{
		Name:  "cmd",
		Usage: "cmd user@ip 'cmd line'",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "t",
				Usage: "--tenv",
			},
			&cli.StringFlag{
				Name:  "p",
				Usage: "--penv",
			},
			&cli.StringFlag{
				Name:  "P",
				Usage: "--port",
			},
		},
		Action: (&util.Cli{}).Run,
	}
}
