package setting

import "time"

// server 配置对象 保存 Port 端口参数
type server struct {
	Port string
}

// session 配置对象 保存 Cookie 名称以及 Cookie 的过期时间
type session struct {
	CookieName   string
	CookieExpire time.Duration
}

// Server 创建服务相关配置
var Server = server{}

// Session 创建会话相关配置
var Session = session{}

// Load 读取配置
func Load() {
	Server.Port = ":3000"
	Session.CookieName = "cafe_echo_example_session_id"
	Session.CookieExpire = time.Hour * 1
}
