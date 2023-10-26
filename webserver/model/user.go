package model

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/labstack/echo"
	"io"
	"io/ioutil"
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

// echo上下文
var e *echo.Echo

// 用于保存数据的map
var users map[ID]User

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

func (u *User) Copy(f *User) {
	u.ID = f.ID
	u.UserID = f.UserID
	u.Password = f.Password
	u.FullName = f.FullName
	u.Roles = make([]Role, len(f.Roles))
	copy(u.Roles, f.Roles)
}

// Start 开启Accessor
func (a *UserDataAccessor) Start(echo *echo.Echo) error {
	e = echo
	users = make(map[ID]User)
	if err := a.decodeJSON(); err != nil {
		return err
	}
	go a.mainLoop()
	return nil
}

// Stop 关闭Accessor对象
func (a *UserDataAccessor) Stop() {
	a.stopCh <- struct{}{}
}

// FindAll 查询所有用户
func (a *UserDataAccessor) FindAll() ([]User, error) {

	respCh := make(chan response, 1)
	defer close(respCh)
	req := []interface{}{}
	cmd := command{commandFindAll, req, respCh}
	a.commandCh <- cmd
	resp := <-respCh
	var res []User
	if resp.err != nil {
		e.Logger.Debugf("User Find Error, [%s]", resp.err)
		return res, resp.err
	}
	if res, ok := resp.result[0].([]User); ok {
		return res, nil
	}
	e.Logger.Debugf("User Find Error [%s]", ErrorOther)
	return nil, ErrorOther
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

// decodeJSON 将JSON数据转换为map
func (a *UserDataAccessor) decodeJSON() error {
	// JSON文件读取
	bytes, err := ioutil.ReadFile("data/users.json")
	if err != nil {
		return err
	}
	// 解析JSON数据
	var recodes []User
	if err := json.Unmarshal(bytes, &recodes); err != nil {
		return err
	}
	// 将解析结果封装为map
	for _, x := range recodes {
		users[x.ID] = x
	}
	return nil
}

// EncodeStringMD5 返回一个MD5封装的字符串
func EncodeStringMD5(str string) StringMD5 {
	h := md5.New()
	io.WriteString(h, str)
	encodeStr := hex.EncodeToString(h.Sum(nil))
	res := StringMD5(encodeStr)
	return res
}

// mainLoop UserDataAccessor的主循环处理
func (a *UserDataAccessor) mainLoop() {
	a.stopCh = make(chan struct{}, 1)
	a.commandCh = make(chan command, 1)
	defer close(a.commandCh)
	defer close(a.stopCh)
	e.Logger.Info("model.UserDataAccessor:start")
loop:
	for {
		// 根据接收的指令进行对应的处理
		select {
		case cmd := <-a.commandCh:
			switch cmd.cmdType {
			// 全部查询
			case commandFindAll:
				results := []User{}
				for _, x := range users {
					user := User{}
					user.Copy(&x)
					results = append(results, user)
				}
				res := []interface{}{results}
				cmd.responseCh <- response{res, nil}
				break
			// ID查询
			case commandFindByID:
				// 未实装
				cmd.responseCh <- response{nil, ErrorNotImplemented}
				break
			// UserID查询
			case commandFindByUserID:
				reqUserID, ok := cmd.req[0].(string)
				if !ok {
					cmd.responseCh <- response{nil, ErrorBadParameter}
					break
				}
				reqOption, ok := cmd.req[1].(FindOption)
				if !ok {
					cmd.responseCh <- response{nil, ErrorBadParameter}
					break
				}
				results := []User{}
				for _, x := range users {
					if x.UserID == reqUserID {
						user := User{}
						user.Copy(&x)
						results = append(results, user)
						if reqOption == FindFirst {
							break
						}
					}
				}
				if len(results) <= 0 {
					cmd.responseCh <- response{nil, ErrorNotFound}
					break
				}
				if reqOption == FindUnique && len(results) > 1 {
					cmd.responseCh <- response{nil, ErrorMultipleResults}
					break
				}
				res := []interface{}{results}
				cmd.responseCh <- response{res, nil}
			default:
				cmd.responseCh <- response{nil, ErrorInvalidCommand}
			}
		case <-a.stopCh:
			break loop
		}
	}
	e.Logger.Info("model.UserDataAccessor:stop")
}
