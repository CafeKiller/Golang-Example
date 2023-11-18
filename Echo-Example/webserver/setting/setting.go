package setting

import "time"

type server struct {
	Port string
}
type session struct {
	CookieName   string
	CookieExpire time.Duration
}

// Server 服务相关配置
var Server = server{}

// Session 会话相关配置
var Session = session{}

// Load 读取配置
func Load() {
	Server.Port = ":3000"
	Session.CookieName = "cafe_echo_example_session_id"
	Session.CookieExpire = time.Hour * 1
}
