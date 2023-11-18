package session

import (
	"errors"
	"github.com/labstack/echo"
	"time"
)

type ID string

// Store 用于保存会话数据与token的结构体
type Store struct {
	Data             map[string]string
	ConsistencyToken string
}

// Manager Session管理对象
type Manager struct {
	stopCh    chan struct{}
	commandCh chan command
	stopGCCh  chan struct{}
}

// 生成Manager返回的实例
var (
	ErrorBadParameter   = errors.New("Bad Parameter")
	ErrorNotFound       = errors.New("Not Found")
	ErrorInvalidToken   = errors.New("Invalid Token")
	ErrorInvalidCommand = errors.New("Invalid Command")
	ErrorNotImplemented = errors.New("Not Implemented")
	ErrorOther          = errors.New("Other")
)

// Start 开启Manager
func (m *Manager) Start(echo *echo.Echo) {
	e = echo
}

// Stop 停止Manager
func (m *Manager) Stop() {
	// 创建一个空对象, 并使用通道将值传递个stopGCCh
	m.stopGCCh <- struct{}{}
	time.Sleep(100 * time.Millisecond)
	// 创建一个空对象, 并使用通道将值传递到stopCh
	m.stopCh <- struct{}{}
}

// Create 创建会话
func (m *Manager) Create() (ID, error) {
	respCh := make(chan response, 1)
	defer close(respCh)
	cmd := command{commandCreate, nil, respCh}
	m.commandCh <- cmd
	resp := <-respCh
	var res ID
	if resp.err != nil {
		e.Logger.Debugf("Session Create Error. [%s]", resp.err)
		return res, resp.err
	}
	if res, ok := resp.result[0].(ID); ok {
		return res, nil
	}
	e.Logger.Debugf("Session Create Error. [%s]", ErrorOther)
	return res, ErrorOther
}

// LoadStore 读取数据存储
func (m *Manager) LoadStore(sessionID ID) (Store, error) {
	respCh := make(chan response, 1)
	defer close(respCh)
	req := []interface{}{sessionID}
	cmd := command{commandLoadStore, req, respCh}
	m.commandCh <- cmd
	resp := <-respCh
	var res Store
	if resp.err != nil {
		e.Logger.Debugf("Session[%s] Load store Error. [%s]", sessionID, resp.err)
		return res, resp.err
	}
	if res, ok := resp.result[0].(Store); ok {
		return res, nil
	}
	e.Logger.Debugf("Session[%s] Load store Error. [%s]", sessionID, ErrorOther)
	return res, ErrorOther
}

// SavaStore 保存数据
func (m *Manager) SavaStore(sessionID ID, sessionStore Store) error {
	respCh := make(chan response, 1)
	defer close(respCh)
	req := []interface{}{sessionID, sessionStore}
	cmd := command{commandSavaStore, req, respCh}
	m.commandCh <- cmd
	resp := <-respCh
	if resp.err != nil {
		e.Logger.Debugf("Session[%s] Sava store Error. [%s]", sessionID, resp.err)
		return resp.err
	}
	return nil
}

// Delete 删除会话
func (m *Manager) Delete(sessionID ID) error {
	respCh := make(chan response, 1)
	defer close(respCh)
	req := []interface{}{sessionID}
	cmd := command{commandDelete, req, respCh}
	m.commandCh <- cmd
	resp := <-respCh
	if resp.err != nil {
		e.Logger.Debugf("Session[%s] Delete Error. [%s]", sessionID, resp.err)
		return resp.err
	}
	return nil
}

// DeleteExpired 删除过期的会话
func (m *Manager) DeleteExpired() error {
	respCh := make(chan response, 1)
	defer close(respCh)
	cmd := command{commandDeleteExpired, nil, respCh}
	m.commandCh <- cmd
	resp := <-respCh
	if resp.err != nil {
		e.Logger.Debugf("Session DeleteExpired Error. [%s]", resp.err)
		return resp.err
	}
	return nil
}
