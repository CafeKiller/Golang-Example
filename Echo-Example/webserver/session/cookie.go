package session

import (
	"Echo-Example/webserver/setting"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

// WriteCookie 将 cookie 写入容器内
func WriteCookie(content echo.Context, sessionID ID) error {

	// 通过 Echo 创建一个 Cookie，
	cookie := new(http.Cookie)
	cookie.Name = setting.Session.CookieName
	cookie.Value = string(sessionID)
	// 设置
	cookie.Expires = time.Now().Add(setting.Session.CookieExpire)
	content.SetCookie(cookie)
	return nil
}

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
