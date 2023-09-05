package main

import (
	"douyin/app"
	"douyin/config"
	"log"
	"os"
)

func main() {
	conf, err := config.Init("config", "yaml")
	if err != nil {
		log.Panicln("初始化配置失败", err)
	}
	config.Print()

	quit := make(chan os.Signal)
	defer close(quit)
	app.Run(conf, quit, nil)
}
