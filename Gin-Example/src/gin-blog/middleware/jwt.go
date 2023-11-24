package middleware

import (
	"Gin-Example/src/gin-blog/pkg/err"
	"Gin-Example/src/gin-blog/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func JWT() gin.HandlerFunc {
	return func(cont *gin.Context) {
		var code int
		var data interface{}

		code = err.SUCCESS
		token := cont.Query("token")
		if token == "" {
			code = err.INVALID_PARAMS
		} else {
			claims, e := util.ParseToken(token)
			if e != nil {
				code = err.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = err.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		if code != err.SUCCESS {
			cont.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  err.GetMsg(code),
				"data": data,
			})
			cont.Abort()
			return
		}

		cont.Next()
	}
}
