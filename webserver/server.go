package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"html/template"
)

var templates map[string]*template.Template

//var userDA *model.UserDataAccessor

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
	echo.Use(middleware.Logger())
	echo.Use(middleware.Recover())

	// 设置静态文件路径
	setStaticRoute(echo)

	//

	fmt.Println("Hello World")

}
