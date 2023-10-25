package model

import (
	"errors"
	"github.com/labstack/echo"
)

// commandType 命令行参数类型
type commandType int

// FindOption 查询参数选项
type FindOption int

// ID 唯一识别信息
type ID string

// StringMD5 通过MD5哈希化后的字符串
type StringMD5 string

// Role 用户权限
type Role string

// response 命令执行结果
type response struct {
	result []interface{}
	err    error
}

const (
	commandFindAll commandType = iota
	commandFindByID
	commandFindByUserID
)

// 搜索时选项
const (
	FIndAll    FindOption = iota // 全件検索
	FindFirst                    // 1件目のみ返す
	FindUnique                   // 結果が1件のみでない場合にはエラーを返す
)

var e *echo.Echo

// DataAccessor返回的异常选项
var (
	ErrorNotFound        = errors.New("Not found")
	ErrorMultipleResults = errors.New("Multiple results")
	ErrorInvalidCommand  = errors.New("Invalid Command")
	ErrorBadParameter    = errors.New("Bad Parameter")
	ErrorNotImplemented  = errors.New("Not Implemented")
	ErrorOther           = errors.New("Other")
)

// 执行命令的参数
type command struct {
	cmdType    commandType
	req        []interface{}
	responseCh chan response
}

// UserDataAccessor 用户数据操作对象
type UserDataAccessor struct {
	stopCh    chan struct{}
	commandCh chan command
}

type User struct {
	ID       ID        `json:"id"`
	UserID   string    `json:"user_id"`
	Password StringMD5 `json:"password"`
	FullName string    `json:"full_name"`
	Roles    []Role    `json:"roles"`
}

// FindByUserID 通过UserID查询对应数据
func (a *UserDataAccessor) FindByUserID(reqUserID string, option FindOption) ([]User, error) {
	respCh := make(chan response, 1)
	defer close(respCh)
	req := []interface{}{reqUserID, option}
	cmd := command{commandFindByID, req, respCh}
	a.commandCh <- cmd
	resp := <-respCh
	var res []User
	if resp.err != nil {
		e.Logger.Debugf("User[User=%s] Find Error. [%s]", reqUserID, resp.err)
		return res, resp.err
	}
	if res, ok := resp.result[0].([]User); ok {
		return res, nil
	}
	e.Logger.Debugf("User[User=%s] Find Error. [%s]", reqUserID, ErrorOther)
	return res, ErrorOther
}
