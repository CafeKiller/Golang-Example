package main

import (
	"Echo-Example/webserver/model"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"html/template"
)

var templates map[string]*template.Template

var userDA *model.UserDataAccessor

func main() {
	// 创建echo对象
	e := echo.New()

	// 设置日志的输入级别
	// e.Logger.SetLevel(log.INFO)
	e.Logger.SetLevel(log.DEBUG)

}
