package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"AI-PROXY/config"
	"AI-PROXY/repository"
	"AI-PROXY/router"
	"AI-PROXY/util"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	//1、解析命令行参数
	configPath := flag.String("config", "config.json", "配置文件路径")
	flag.Parse()

	//2、加载配置
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败：%v", err)
	}

	//3、初始化日志
	util.InitLogger()

	//4、连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Database,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		util.Logger.Fatalf("连接数据库失败：%v", err)
	}

	//5、初始化repository层
	repository.InitDB(db)

	//6.初始化路由
	r := router.SetupRouter()

	//7、启动HTTP服务
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	go func() {
		util.Logger.Infof("HTTP服务器启动,监听地址: %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			util.Logger.Fatalf("HTTP服务器启动失败: %v", err)
		}
	}()

	// 8. 等待中断信号，优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	util.Logger.Info("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		util.Logger.Fatalf("服务器关闭失败: %v", err)
	}
	util.Logger.Info("服务器已关闭")
}
