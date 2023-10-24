package session

import "errors"

type ID string

// Store 用于保存会话数据与token的结构体
type Store struct {
	Data             map[string]string
	ConsistencyToken string
}

// Manager Session管理对象
type Manager struct {
	stopCh    chan struct{}
	commandCH chan command
	stopGCCh  chan struct{}
}

// 生成Manager返回的实例
var (
	ErrorOther = errors.New("Other")
)

// LoadStore 读取数据存储
func (m *Manager) LoadStore(sessionID ID) (Store, error) {
	respCh := make(chan response, 1)
	defer close(respCh)
	req := []interface{}{sessionID}
	cmd := command{commandLoadStore, req, respCh}
	m.commandCH <- cmd
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
