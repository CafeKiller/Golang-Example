package session

import (
	"Echo-Example/webserver/setting"
	"github.com/labstack/echo"
)

// ReadCookie 从echo容器内读取cookie
func ReadCookie(ctx echo.Context) (ID, error) {

	var sessionID ID
	cookie, err := ctx.Cookie(setting.Session.CookieName)
	if err != nil {
		return sessionID, err
	}
	sessionID = ID(cookie.Value)
	return sessionID, nil

}
