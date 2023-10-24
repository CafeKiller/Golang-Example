package main

import (
	"Echo-Example/webserver/session"
	"errors"
	"github.com/labstack/echo"
)

var (
	ErrorNotLoggedIn   = errors.New("Not Logged In")
	ErrorInvalidUserID = errors.New("Invalid UserID")
)

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
