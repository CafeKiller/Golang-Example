package session

import (
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	"time"
)

var e *echo.Echo

const (
	commandCreate        commandType = iota // 创建会话
	commandLoadStore                        // 读取会话数据
	commandSavaStore                        // 保存数据
	commandDelete                           // 删除数据
	commandDeleteExpired                    // 数据有效时间
)

var e *echo.Echo

// 命令类型指定
type commandType int

// 会话信息
type session struct {
	store  Store
	expire time.Time
}

// 会话的有效时间
const sessionExpire = 3 * time.Minute

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

// mainLoop sessionManager主循环处理
func (m *Manager) mainLoop() {

	session := make(map[ID]session)
	m.stopCh = make(chan struct{}, 1)
	m.commandCh = make(chan command, 1)
	defer close(m.commandCh)
	defer close(m.stopCh)
	e.Logger.Info("session.Manager:start")
loop:
	for {
		select {
		case cmd := <-m.commandCh:
			switch cmd.cmdType {
			// 创建会话数据
			case commandCreate:

			// 读取会话数据
			case commandLoadStore:

			// 保存会话数据
			case commandSavaStore:

			// 删除会话数据
			case commandDelete:

			// 删除到期的会话数据
			case commandDeleteExpired:

			// 默认行为
			default:
				cmd.responseCh <- response{nil, ErrorInvalidCommand}
			}
		case <-m.stopCh:
			break loop
		}

	}
	e.Logger.Info("session.Manager:stop")
}

// gcLoop 会话到期时,定时进行处理
func (m *Manager) gcLoop() {
	m.stopGCCh = make(chan struct{}, 1)
	defer close(m.stopGCCh)
	e.Logger.Info("session.Manager GC:start")
	t := time.NewTicker(1 * time.Minute)
loop:
	for {
		select {
		case <-t.C:
			respCh := make(chan response, 1)
			defer close(respCh)
			cmd := command{commandDeleteExpired, nil, respCh}
			m.commandCh <- cmd
			<-respCh
		case <-m.stopGCCh:
			break loop
		}
	}
	t.Stop()
	e.Logger.Info("session.Manager GC:stop")
}

// createSessionID 创建新的会话ID
func createSessionID() string {
	return uuid.NewV4().String()
}

// createToken 创建新的token
func createToken() string {
	return uuid.NewV4().String()
}
