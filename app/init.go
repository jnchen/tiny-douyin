package app

import (
	"context"
	"douyin/config"
	"douyin/db"
	"douyin/router"
	"douyin/service"
	"douyin/storage"
	"errors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(
	conf *config.TotalConfig,
	quit chan os.Signal,
	running *chan struct{},
) {
	var err error

	err = db.Init(conf.MySQL)
	if err != nil {
		log.Panicln("初始化数据库失败", err)
	}
	defer func() {
		if err = db.SQL().Close(); err != nil {
			log.Println("关闭数据库连接失败", err)
		}
		log.Println("关闭数据库连接")
	}()

	err = storage.Init(conf.Storage)
	if err != nil {
		log.Panicln("初始化存储失败", err)
	}

	quitMsgSrv := make(chan struct{})
	defer close(quitMsgSrv)
	go service.RunMessageServer(quitMsgSrv)

	if conf.Release {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	if !conf.Release {
		pprof.Register(r)
	}
	err = router.Init(r)
	if err != nil {
		log.Panicln("初始化路由失败", err)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	defer func() {
		log.Println("关闭服务中……")
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		if err = srv.Shutdown(ctx); err != nil {
			log.Fatal("服务关闭失败：", err)
		}
		<-ctx.Done()
	}()
	go func() {
		if err = srv.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
			log.Println("服务关闭")
		} else if err != nil {
			log.Fatalf("启动服务失败：%v", err)
		}
	}()

	// 等待中断信号后关闭服务器，同时为关闭服务器操作设置一个超时时间
	signal.Notify(
		quit,
		syscall.SIGINT,  // kill -2 是 syscall.SIGINT
		syscall.SIGTERM, // kill（无参数）默认发送 syscall.SIGTERM
		// kill -9 是 syscall.SIGKILL，但无法捕获，因此不需要添加它
	)
	if running != nil {
		*running <- struct{}{} // 通知服务已启动
		close(*running)
	}
	<-quit
	quitMsgSrv <- struct{}{} // 通知消息服务器关闭
}
