package api

import (
	"Gin-Example/src/gin-blog/models"
	"Gin-Example/src/gin-blog/pkg/err"
	"Gin-Example/src/gin-blog/pkg/util"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

// GetAuth
func GetAuth(cont *gin.Context) {
	username := cont.Query("username")
	password := cont.Query("password")

	valid := validation.Validation{}
	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	data := make(map[string]interface{})
	code := err.INVALID_PARAMS
	if ok {
		isExist := models.CheckAuth(username, password)
		if isExist {
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
		for _, err := range valid.Errors {
			logs.Info(err.Key, err.Message)
		}
	}

	cont.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  err.GetMsg(code),
		"data": data,
	})

}
