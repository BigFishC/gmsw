package main

import (
	"fmt"

	"github.com/liuchong/gmsw/config"
	"github.com/liuchong/gmsw/secret"
	"github.com/urfave/cli/v2"
)

func main() {

	app := cli.NewApp()
	app.Name = "gmsw"
	app.Usage = "forgeted"
	app.Commands = []*cli.Command{}
	config.LoadConfig(".")

	decrypt, _ := secret.DecryptByAes(config.ConfigStruct.Testpwd, secret.PwdKey)

	fmt.Printf("解密后：%s\n", decrypt)
}
