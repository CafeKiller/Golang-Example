package session

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

// LoadStore 读取数据存储
func (m *Manager) LoadStore(sessionID ID) (Store, error) {
	respCh := make(chan response, 1)
	defer close(respCh)
	req := []interface{}{sessionID}
	cmd := command{}
}
