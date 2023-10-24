package session

const {
	commandLoadStore = 1	// 读取数据
}

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
