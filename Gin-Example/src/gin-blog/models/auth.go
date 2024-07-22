package models

type Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// CheckAuth 身份检验
func CheckAuth(username, password string) bool {
	// 定义Auth结构体
	var auth Auth
	// 从数据库中查询id，where条件为用户名和密码
	db.Select("id").Where(Auth{Username: username, Password: password}).First(&auth)
	// 如果id大于0，则返回true
	if auth.ID > 0 {
		return true
	}

	return false
}
