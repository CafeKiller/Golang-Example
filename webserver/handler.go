package main

import (
	"github.com/labstack/echo"
	"net/http"
)

// setRoute 设置路由处理程序
func setRoute(echo *echo.Echo) {
	echo.GET("/", handleIndexGet)
}

// GET :/
func handleIndexGet(content echo.Context) error {
	return content.Render(http.StatusOK, "index", "world")
}

// GET :/users/:user_id
// POST :/users/:user_id
func handleUsers(content echo.Context) error {
	userID := content.Param("user_id")
	err := CheckUserID(content, userID)
	if err != nil {
		content.Echo().Logger.Debugf("User page[%s] Role Error. [%s]", userID, err)
		mes := "未登录 / Not logged in / ログインしていません。"
		return content.Render(http.StatusOK, "error", mes)
	}
}
