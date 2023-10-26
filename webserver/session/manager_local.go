package session

import "github.com/labstack/echo"

var e *echo.Echo

const (
	commandCreate        commandType = iota // 创建会话
	commandLoadStore                        // 读取会话数据
	commandSavaStore                        // 保存数据
	commandDelete                           // 删除数据
	commandDeleteExpired                    // 数据有效时间
)

// 命令类型指定
type commandType int

// 执行命令参数
type command struct {
	cmdType    commandType
	req        []interface{}
	responseCh chan response
}

// 命令执行结构
type response struct {
	result []interface{}
	err    error
}
