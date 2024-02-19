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

// 声明一个全局模版集合 ( map 集合, 一种无序键值对集合)
var templates map[string]*template.Template

// 创建全局 session 管理器
var sessionManager *session.Manager

var userDA *model.UserDataAccessor

func main() {
	// 创建 echo 对象
	echo := echo.New()

	// 设置 echo 的内置日志管理器的输入级别
	// e.Logger.SetLevel(log.INFO)
	echo.Logger.SetLevel(log.DEBUG)

	// 实例化全局模版集合
	t := &Template{}
	// 将模版实例绑定至 echo 的 Renderer 接口上, echo 会自行调用渲染
	echo.Renderer = t

	// 设置中间件
	// Logger 中间件记录有关每个 HTTP 请求的信息。
	echo.Use(middleware.Logger())
	// Recover 中间件从 panic 链中的任意位置恢复程序， 打印堆栈的错误信息，并将错误集中交给 HTTPErrorHandler 处理。
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
