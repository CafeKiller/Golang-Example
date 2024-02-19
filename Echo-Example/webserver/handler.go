package main

import (
	"Echo-Example/webserver/model"
	"github.com/labstack/echo"
	"net/http"
)

// setRoute 传入 echo 对象, 并为其绑定对应路由处理程序
func setRoute(echo *echo.Echo) {

	echo.GET("/", handleIndexGet)
	echo.GET("/login", handleLoginGet)
	echo.POST("/login", handleLoginPost)
	echo.GET("/users/:user_id", handleUsers)
	echo.POST("/users/:user_id", handleUsers)

	// 创建分组 Group, 设置为管理员专属访问页面
	admin := echo.Group("/admin", MiddlewareAuthAdmin)
	admin.GET("", handleAdmin)
	admin.POST("", handleAdmin)
	admin.GET("/users", handleAdminUsersGet)
}

// GET :/
// handleIndexGet 用于处理主页 index 的
func handleIndexGet(content echo.Context) error {
	return content.Render(http.StatusOK, "index", "world")
}

// GET :/users/:user_id
// POST :/users/:user_id
// handleUsers 用于处理用户页 users 的数据
func handleUsers(content echo.Context) error {
	// 获取参数 user_id
	userID := content.Param("user_id")
	// 校验 user_id 的合法性
	err := CheckUserID(content, userID)
	if err != nil {
		content.Echo().Logger.Debugf("User page[%s] Role Error. [%s]", userID, err)
		mes := "未登录 / Not logged in / ログインしていません。"
		return content.Render(http.StatusOK, "error", mes)
	}
	users, err := userDA.FindByUserID(content.Param("user_id"), model.FindFirst)
	if err != nil {
		return content.Render(http.StatusOK, "error", err)
	}
	user := users[0]
	return content.Render(http.StatusOK, "user", user)
}

// GET:/admin
// POST:/admin
// handleAdmin 用于处理管理员主页 admin 的数据
func handleAdmin(content echo.Context) error {
	return content.Render(http.StatusOK, "admin", nil)
}

// GET:/admin/users
// handleAdminUsersGet 用于处理管理端用户页 admin/users 的数据
func handleAdminUsersGet(content echo.Context) error {
	// 查询所有用户的数据
	users, err := userDA.FindAll()
	if err != nil {
		return err
	}
	return content.Render(http.StatusOK, "admin_users", users)
}

// GET:/login
// handleLoginGet 用于处理登录页 login 的数据
func handleLoginGet(content echo.Context) error {
	return content.Render(http.StatusOK, "login", nil)
}

// POST:/login
func handleLoginPost(content echo.Context) error {
	userID := content.FormValue("userid")
	password := content.FormValue("password")
	err := UserLogin(content, userID, password)
	if err != nil {
		content.Echo().Logger.Debugf("User[%s] Login Error. [%s]", userID, err)
		msg := "用户ID或密码错误"
		data := map[string]string{"user_id": userID, "password": password, "msg": msg}
		return content.Render(http.StatusOK, "login", data)
	}
	// 检查登录的用户是否为管理员
	isAdmin, err := CheckRoleByUserID(userID, model.RoleAdmin)
	if err != nil {
		content.Echo().Logger.Debugf("Admin Role Check Error. [%s]", userID, err)
		isAdmin = false
	}
	if isAdmin {
		// 使用管理员登录时, 跳转至管理员页面
		content.Echo().Logger.Debugf("User is Admin [%s]", userID)
		return content.Redirect(http.StatusTemporaryRedirect, "/admin")
	}
	return content.Redirect(http.StatusTemporaryRedirect, "/users/"+userID)
}

// POST:/logout
func handleLogoutPost(content echo.Context) error {
	err := UserLogout(content)
	if err != nil {
		content.Echo().Logger.Debugf("User Logout Error. [%s]", err)
		return content.Render(http.StatusOK, "login", nil)
	}
	msg := "退出登录"
	data := map[string]string{"user_id": "", "password": "", "msg": msg}
	return content.Render(http.StatusOK, "login", data)
}
