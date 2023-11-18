package main

import (
	"Echo-Example/webserver/model"
	"Echo-Example/webserver/session"
	"errors"
	"github.com/labstack/echo"
	"net/http"
)

var (
	ErrorNotLoggedIn     = errors.New("Not Logged In")
	ErrorInvalidUserID   = errors.New("Invalid UserID")
	ErrorInvalidPassword = errors.New("Invalid Password")
)

// UserLogin 用户登录处理
func UserLogin(content echo.Context, userID string, password string) error {

	users, err := userDA.FindByUserID(userID, model.FindFirst)
	if err != nil {
		return err
	}
	user := &users[0]
	encodePassword := model.EncodeStringMD5(password)
	if user.Password != encodePassword {
		return ErrorInvalidPassword
	}
	sessionID, err := sessionManager.Create()
	if err != nil {
		return err
	}
	err = session.WriteCookie(content, sessionID)
	if err != nil {
		return err
	}
	sessionStore, err := sessionManager.LoadStore(sessionID)
	if err != nil {
		return err
	}
	sessionData := map[string]string{
		"user_id": userID,
	}
	sessionStore.Data = sessionData
	err = sessionManager.SavaStore(sessionID, sessionStore)
	if err != nil {
		return err
	}
	return nil
}

// UserLogout 用户退出
func UserLogout(content echo.Context) error {
	sessionID, err := session.ReadCookie(content)
	if err != nil {
		return err
	}
	err = sessionManager.Delete(sessionID)
	if err != nil {
		return err
	}
	return nil
}

// CheckUserID 判断UserID是否合法
func CheckUserID(ctx echo.Context, userID string) error {
	sessionID, err := session.ReadCookie(ctx)
	if err != nil {
		return err
	}
	sessionStore, err := sessionManager.LoadStore(sessionID)
	if err != nil {
		return err
	}
	sessionUserID, ok := sessionStore.Data["user_id"]
	if !ok {
		return ErrorNotLoggedIn
	}
	if sessionUserID != userID {
		return ErrorInvalidUserID
	}
	return nil
}

// CheckRole 检查用户是否拥有登录权限
func CheckRole(content echo.Context, role model.Role) (bool, error) {
	sessionID, err := session.ReadCookie(content)
	if err != nil {
		return false, err
	}
	sessionStore, err := sessionManager.LoadStore(sessionID)
	if err != nil {
		return false, err
	}
	sessionUserID, ok := sessionStore.Data["user_id"]
	if !ok {
		return false, ErrorNotLoggedIn
	}
	havaRole, err := CheckRoleByUserID(sessionUserID, role)
	return havaRole, nil
}

// CheckRoleByUserID 确认指定用户是否拥有权限
func CheckRoleByUserID(userID string, role model.Role) (bool, error) {
	users, err := userDA.FindByUserID(userID, model.FindFirst)
	if err != nil {
		return false, err
	}
	user := &users[0]
	for _, v := range user.Roles {
		if v == role {
			return true, nil
		}
	}
	return false, nil
}

// MiddlewareAuthAdmin 只允许管理员登录的页面
func MiddlewareAuthAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		isAdmin, err := CheckRole(context, model.RoleAdmin)
		if err != nil {
			context.Echo().Logger.Debugf("Admin Page Role Error. [%s]", err)
			isAdmin = false
		}
		if !isAdmin {
			msg := "权限不足,请使用管理员登录"
			return context.Render(http.StatusOK, "error", msg)
		}
		return next(context)
	}
}
