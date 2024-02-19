package main

import "github.com/labstack/echo"

// setStaticRoute 设置静态文件路径
func setStaticRoute(echo *echo.Echo) {
	// 使用 echo 内置的 Static 接口, 设置各个静态文件的映射路径
	echo.Static("/public/css/", "./public/css/")
	echo.Static("/public/js/", "./public/js/")
	echo.Static("/public/img/", "./public/img/")
}
