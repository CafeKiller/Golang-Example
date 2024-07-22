package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

var (
	Cfg *ini.File

	RunMode string

	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	PageSize  int
	JwtSecret string
)

// init 初始化
func init() {
	// 加载配置文件
	var err error
	Cfg, err = ini.Load("./src/gin-blog/conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini' : %v", err)
	}

	// 加载基本配置
	LoadBase()
	// 加载服务器配置
	LoadServer()
	// 加载应用配置
	LoadApp()
}

// LoadBase 加载基本配置
func LoadBase() {

	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")

}

// LoadServer 加载服务配置
func LoadServer() {

	// 获取 "server" 配置节
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server' : %v", err)
	}

	// 设置RunMode
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")

	// 设置HTTPPort
	HTTPPort = sec.Key("HTTP_PORT").MustInt(8000)
	// 设置ReadTimeout
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	// 设置WriteTimeout
	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}

// LoadApp 加载应用配置
func LoadApp() {

	// 获取 "app" 配置节
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app' : %v", err)
	}

	// 获取JwtSecret
	JwtSecret = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")
	// 获取PageSize
	PageSize = sec.Key("PAGE_SIZE").MustInt(10)
}
