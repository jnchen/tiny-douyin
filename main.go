package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stickit/douyin/router"
	"github.com/stickit/douyin/service"
)

func main() {
	go service.RunMessageServer()

	r := gin.Default()

	router.InitRouter(r)

	err := r.Run()
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
