package config

import (
	"fmt"
	"github.com/spf13/viper" //配置管理
)

var Conf = new(TotalConfig)

type TotalConfig struct {
	*MySQLConfig `mapstructure:"mysql"`
}

type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Dbname   string `mapstructure:"dbname"`
	Options  string `mapstructure:"options"`
	Port     int    `mapstructure:"port"`
}

func init() {
	fmt.Println("解析配置")
	viper.SetConfigFile("config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("配置读取失败, err:%v\n", err)
		return
	}

	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("配置解析失败, err:%v\n", err)
	}
	fmt.Println(Conf.MySQLConfig)
}
