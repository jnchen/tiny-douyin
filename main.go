package main

import (
	"douyin/db"
	"douyin/router"
	"douyin/service"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	defer func() {
		if err := db.SQL().Close(); err != nil {
			log.Println("关闭数据库连接失败", err)
		}
		log.Println("关闭数据库连接")
	}()

	go service.RunMessageServer()

	r := gin.Default()
	pprof.Register(r)
	router.InitRouter(r)

	err := r.Run(":8080")
	if err != nil {
		log.Panicln("启动服务失败", err)
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
