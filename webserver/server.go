package main

import (
	"Echo-Example/webserver/model"
	"Echo-Example/webserver/session"
	"Echo-Example/webserver/setting"
	"context"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"html/template"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 声明一个全局模版对象
var templates map[string]*template.Template

// 全局session管理器
var sessionManager *session.Manager

var userDA *model.UserDataAccessor

func main() {
	// 创建echo对象
	echo := echo.New()

	// 设置日志的输入级别
	// e.Logger.SetLevel(log.INFO)
	echo.Logger.SetLevel(log.DEBUG)

	// 使用echo内置的模版渲染
	t := &Template{}
	echo.Renderer = t

	// 设置中间件
	// 中间件 用于记录每一个HTTP请求信息
	echo.Use(middleware.Logger())
	// 中间件 用于从panic错误链中恢复程序,打印错误信息,并将错误集中到 HTTPErrorHandler 处理
	echo.Use(middleware.Recover())

	// 设置静态文件路径
	setStaticRoute(echo)

	// 设置各路由的处理程序
	setRoute(echo)

	// 开启会话管理
	sessionManager = &session.Manager{}
	sessionManager.Start(echo)

	// 开启数据访问对象
	userDA = &model.UserDataAccessor{}
	userDA.Start(echo)

	// 启动服务器
	go func() {
		if err := echo.Start(setting.Server.Port); err != nil {
			echo.Logger.Info("shutting down the server")
		}
	}()

	// 检测到中断时, 10秒内无响应就自动关闭服务
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := echo.Shutdown(ctx); err != nil {
		echo.Logger.Info(err)
		echo.Close()
	}

	// 关闭数据访问对象
	userDA.Stop()

	// 关闭会话管理
	sessionManager.Stop()

	// 在结束后稍等停止等待一下
	time.Sleep(1 * time.Second)
}

// init 初始化函数
func init() {
	// 加载配置
	setting.Load()
	// 加载HTML模版
	loadTemplates()
}
