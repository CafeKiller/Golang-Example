package main

import (
	"Echo-Example/webserver/model"
	"Echo-Example/webserver/session"
	"errors"
	"github.com/labstack/echo"
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
