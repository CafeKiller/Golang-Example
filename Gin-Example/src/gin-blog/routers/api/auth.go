package api

import (
	"Gin-Example/src/gin-blog/models"
	"Gin-Example/src/gin-blog/pkg/err"
	"Gin-Example/src/gin-blog/pkg/logging"
	"Gin-Example/src/gin-blog/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

// GetAuth 获取token
func GetAuth(cont *gin.Context) {

	// 获取请求参数
	username := cont.Query("username")
	password := cont.Query("password")

	// 验证参数
	valid := validation.Validation{}
	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	// 初始化返回数据
	data := make(map[string]interface{})
	code := err.INVALID_PARAMS

	// 验证通过
	if ok {
		// 检查用户名和密码是否正确
		isExist := models.CheckAuth(username, password)
		if isExist {
			// 生成token
			token, e := util.GenerateToken(username, password)
			if e != nil {
				code = err.ERROR_AUTH_TOKEN
			} else {
				data["token"] = token
				code = err.SUCCESS
			}
		} else {
			code = err.ERROR_AUTH
		}
	} else {
		// 验证失败
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	// 返回结果
	cont.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  err.GetMsg(code),
		"data": data,
	})
}
