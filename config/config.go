package config

import (
	"fmt"
	"log"

	"github.com/BigFishC/gmsw/secret"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

// ConfigStruct 配置文件内容
type ConfigStruct struct {
	Tpwd string `yaml:"tpwd"`
	Ppwd string `yaml:"ppwd"`
}

// LoadConfig 加载配置文件
func (c *ConfigStruct) LoadConfig() error {
	configVip := OptConfig("/bitnami/jenkins/conf.yml")
	if err := configVip.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := configVip.Unmarshal(c); err != nil {
		panic(err)
	}
	return nil
}

// OptConfig 操作配置文件
func OptConfig(confName string) *viper.Viper {
	configVip := viper.New()
	configVip.SetConfigFile(confName)
	return configVip
}

// WriteEncryptPwd 密码加密写入配置方法
func WriteEncryptPwd(param string, pwd string) error {
	if pwd == "" {
		fmt.Printf("%s is nil", param)
	} else {
		configVip := OptConfig("/bitnami/jenkins/conf.yml")
		if err2 := configVip.ReadInConfig(); err2 != nil {
			if _, ok := err2.(viper.ConfigFileNotFoundError); ok {
				panic("配置文件未找到！")
			} else {
				panic("找到配置文件,但是解析错误！")
			}
		}
		configVip.Set(param, pwd)
		configVip.WriteConfig()
		log.Fatal("Encrypt password successful!")
	}
	return nil
}

func WriteEnvChange(param string, cli *cli.Context) error {
	password := cli.String(param)
	origin := []byte(password)
	encrypt, _ := secret.EncryptByAes(origin, secret.PwdKey)
	WriteEncryptPwd(param, encrypt)
	return nil
}

// UpdateConfig 更新配置文件
func (c *ConfigStruct) UpdateConfig(cli *cli.Context) error {

	if cli.NArg() == 0 && cli.NumFlags() == 1 {
		env := cli.FlagNames()[0]
		switch env {
		case "tpwd":
			WriteEnvChange("tpwd", cli)
		case "ppwd":
			WriteEnvChange("ppwd", cli)
		default:
			log.Fatal("Please use the -h parameter for help")
		}
	} else {
		log.Fatal("Please use the -h parameter for help")
	}
	return nil
}
