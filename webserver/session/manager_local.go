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

	sessions := make(map[ID]session)
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
				sessionID := ID(createSessionID())
				session := session{}
				sessionStore := Store{}
				sessionData := make(map[string]string)
				sessionStore.Data = sessionData
				sessionStore.ConsistencyToken = createToken()
				session.store = sessionStore
				session.expire = time.Now().Add(sessionExpire)
				sessions[sessionID] = session
				res := []interface{}{sessionID}
				e.Logger.Debugf("Session[%s] Create. expire[%s]", sessionID, session.expire)
				cmd.responseCh <- response{res, nil}
			// 读取会话数据
			case commandLoadStore:
				reqSessionID, ok := cmd.req[0].(ID)
				if !ok {
					cmd.responseCh <- response{nil, ErrorBadParameter}
					break
				}
				session, ok := sessions[reqSessionID]
				if !ok {
					cmd.responseCh <- response{nil, ErrorNotFound}
					break
				}
				if time.Now().After(session.expire) {
					cmd.responseCh <- response{nil, ErrorNotFound}
					break
				}
				sessionStore := Store{}
				sessionData := make(map[string]string)
				for k, v := range session.store.Data {
					sessionData[k] = v
				}
				sessionStore.Data = sessionData
				sessionStore.ConsistencyToken = session.store.ConsistencyToken
				session.expire = time.Now().Add(sessionExpire)
				sessions[reqSessionID] = session
				e.Logger.Debugf("Session[%s] Load store. store[%s] expire[%s]", reqSessionID, session.store, session.expire)
				res := []interface{}{sessionStore}
				cmd.responseCh <- response{res, nil}
			// 保存会话数据
			case commandSavaStore:
				reqSessionID, ok := cmd.req[0].(ID)
				if !ok {
					cmd.responseCh <- response{nil, ErrorBadParameter}
					break
				}
				reqSessionStore, ok := cmd.req[1].(Store)
				if !ok {
					cmd.responseCh <- response{nil, ErrorBadParameter}
					break
				}
				session, ok := sessions[reqSessionID]
				if !ok {
					cmd.responseCh <- response{nil, ErrorNotFound}
					break
				}
				if time.Now().After(session.expire) {
					cmd.responseCh <- response{nil, ErrorNotFound}
					break
				}
				if session.store.ConsistencyToken != reqSessionStore.ConsistencyToken {
					cmd.responseCh <- response{nil, ErrorInvalidToken}
					break
				}
				sessionStore := Store{}
				sessionData := make(map[string]string)
				for k, v := range reqSessionStore.Data {
					sessionData[k] = v
				}
				sessionStore.Data = sessionData
				sessionStore.ConsistencyToken = createToken()
				session.store = sessionStore
				session.expire = time.Now().Add(sessionExpire)
				sessions[reqSessionID] = session
				e.Logger.Debugf("Session[%s] Save store. store[%s] expire[%s]", reqSessionID, session.store, session.expire)
				cmd.responseCh <- response{nil, nil}
			// 删除会话数据
			case commandDelete:
				reqSessionID, ok := cmd.req[0].(ID)
				if !ok {
					cmd.responseCh <- response{nil, ErrorBadParameter}
					break
				}
				session, ok := sessions[reqSessionID]
				if !ok {
					cmd.responseCh <- response{nil, ErrorNotFound}
					break
				}
				if time.Now().After(session.expire) {
					cmd.responseCh <- response{nil, ErrorNotFound}
					break
				}
				delete(sessions, reqSessionID)
				e.Logger.Debugf("Session[%s] Delete.", reqSessionID)
				cmd.responseCh <- response{nil, nil}
			// 删除到期的会话数据
			case commandDeleteExpired:
				e.Logger.Debugf("Run Session GC. Now[%s]", time.Now())
				for k, v := range sessions {
					if time.Now().After(v.expire) {
						e.Logger.Debugf("Session[%s] expire delete. expire[%s]", k, v.expire)
						delete(sessions, k)
					}
				}
				cmd.responseCh <- response{nil, nil}
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
