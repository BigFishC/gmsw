package main

import (
	"log"
	"os"

	"github.com/BigFishC/gmsw/config"
	"github.com/BigFishC/gmsw/proc"
	"github.com/BigFishC/gmsw/service"
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
		KillProcess(),
		StartService(),
	}

	app.Run(os.Args)

}

func Encrypt() *cli.Command {
	return &cli.Command{
		Name:      "encrypt",
		Usage:     "Encrypt the string to conf.yml",
		UsageText: "gmsf encrypt --tpwd=string | --ppwd=string",
		Action:    (&config.ConfigStruct{}).UpdateConfig,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "tpwd",
				Usage: "--tpwd=`STRING`",
			},
			&cli.StringFlag{
				Name:  "ppwd",
				Usage: "--ppwd=`STRING`",
			},
		},
	}
}

func RunCmd() *cli.Command {
	return &cli.Command{
		Name:      "cmd",
		Usage:     "Run commands remotely and transfer files to a remote computer",
		UsageText: "gmsf cmd [-P] [-T] [-t | -p] user@ip 'something'",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "t",
				Usage: "--tenv` `",
			},
			&cli.StringFlag{
				Name:  "p",
				Usage: "--penv` `",
			},
			&cli.StringFlag{
				Name:  "P",
				Usage: "--port=`PORT`",
			},
			&cli.StringFlag{
				Name:  "T",
				Usage: "--trans=`FILENAME`",
			},
		},
		Action: (&util.Cli{}).Server,
	}
}

func KillProcess() *cli.Command {
	return &cli.Command{
		Name:      "kill",
		Usage:     "Kill servicename",
		UsageText: "gmsf kill servicename",
		Action: func(c *cli.Context) error {
			if c.NArg() > 0 {
				pname := c.Args().First()
				proc.KillProcess(pname)
			} else {
				log.Fatal("Please use the -h parameter for help")
			}

			return nil
		},
	}
}

func StartService() *cli.Command {
	return &cli.Command{
		Name:      "start",
		Usage:     "A command in the specified directory",
		UsageText: "gmsf start -d directory -c cmd",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "d",
				Usage: "run `directory`",
			},
			&cli.StringFlag{
				Name:  "c",
				Usage: "run `cmdline`",
			},
		},
		Action: (&service.SShell{}).StartUp,
	}
}
