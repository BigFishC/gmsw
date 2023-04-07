package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

//ConfigStruct 配置文件内容
var ConfigStruct = struct {
	Testpwd string `json:"testpwd"`
	Propwd  string `json:"propwd"`
}{}

//LoadConfig 加载配置文件
func LoadConfig(configpath string) error {
	configVip := viper.New()
	configVip.AddConfigPath(configpath)
	configVip.SetConfigFile("conf.yml")
	if err := configVip.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := configVip.Unmarshal(&ConfigStruct); err != nil {
		panic(err)
	}
	log := zap.S()
	log.Info("The currunt configuration  file is: ")
	log.Infof("%+v", ConfigStruct)
	return nil
}
