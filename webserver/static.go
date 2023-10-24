package main

import "github.com/labstack/echo"

// 设置静态文件路径
func setStaticRoute(echo *echo.Echo) {
	echo.Static("/public/css/", "./public/css/")
	echo.Static("/public/js/", "./public/js/")
	echo.Static("/public/img/", "./public/img/")
}
