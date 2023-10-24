package main

import (
	"Echo-Example/webserver/session"
	"github.com/labstack/echo"
)

// CheckUserID 判断UserID是否合法
func CheckUserID(ctx echo.Context, userID string) error {
	sessionID, err := session.ReadCookie(ctx)
	if err != nil {
		return err
	}

}
