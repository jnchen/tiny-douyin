package main

import (
	"douyin/config"
	"douyin/model"
	"douyin/router"
	"douyin/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化Database
	model.InitDatabase(config.Conf.MySQLConfig)

	go service.RunMessageServer()

	r := gin.Default()

	router.InitRouter(r)

	err := r.Run()
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
