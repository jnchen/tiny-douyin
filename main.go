package main

import (
	"douyin/db"
	"douyin/router"
	"douyin/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		log.Println("启动pprof服务")
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()

	defer func() {
		if err := db.SQL().Close(); err != nil {
			log.Println("关闭数据库连接失败", err)
		}
		log.Println("关闭数据库连接")
	}()

	go service.RunMessageServer()

	r := gin.Default()

	router.InitRouter(r)

	err := r.Run()
	if err != nil {
		log.Panicln("启动服务失败", err)
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
