package main

import (
	"os"

	"github.com/BigFishC/gmsw/config"
	"github.com/BigFishC/gmsw/proc"
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
	}

	app.Run(os.Args)

}

func Encrypt() *cli.Command {
	return &cli.Command{
		Name:   "encrypt",
		Usage:  "encrypt --tpwd=string | --ppwd=string",
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
		Usage: "cmd -P  -T  -t | -p user@ip 'something'",
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
			&cli.StringFlag{
				Name:  "T",
				Usage: "--trans",
			},
		},
		Action: (&util.Cli{}).Server,
	}
}

func KillProcess() *cli.Command {
	return &cli.Command{
		Name:  "kill",
		Usage: "kill servicename",
		Action: func(c *cli.Context) error {
			pname := c.Args().First()
			proc.KillProcess(pname)
			return nil
		},
		// Flags: []cli.Flag{
		// 	&cli.StringFlag{
		// 		Name:  "tpwd",
		// 		Usage: "--tpwd",
		// 	},
		// 	&cli.StringFlag{
		// 		Name:  "ppwd",
		// 		Usage: "--ppwd",
		// 	},
		// },
	}
}
