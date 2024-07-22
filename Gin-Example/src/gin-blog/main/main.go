package main

import (
	"Gin-Example/src/gin-blog/pkg/setting"
	"Gin-Example/src/gin-blog/routers"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
)

// 定义主函数
func main() {

	// 初始化路由器
	router := routers.InitRouter()

	// 创建HTTP服务器
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// 启动服务器
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("Listen: %s \n", err)
		}
	}()

	// 接收系统信号，优雅关闭服务器
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ......")

	// 等待5秒关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// server.Shutdown函数会关闭服务器，但不会立即退出，而是等待所有正在处理的请求完成。
	// 这使得在服务器关闭之前能够完成所有未完成的请求，从而避免潜在的错误。
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}

// GetAppPath 获取应用程序的路径
func GetAppPath() string {
	// 使用 exec.LookPath 获取命令行参数第一个参数的绝对路径
	file, _ := exec.LookPath(os.Args[0])
	// 使用 filepath.Abs 获取文件的绝对路径
	path, _ := filepath.Abs(file)
	// 获取路径中最后一个分隔符的位置
	index := strings.LastIndex(path, string(os.PathSeparator))
	// 返回路径中最后一个分隔符之前的所有字符
	return path[:index]
}
