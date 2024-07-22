package jwt

import (
	"Gin-Example/src/gin-blog/pkg/err"
	"Gin-Example/src/gin-blog/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// JWT 是一个 gin.HandlerFunc 函数，用于检查查询参数中的 JWT 标记
func JWT() gin.HandlerFunc {
	return func(cont *gin.Context) {
		var code int
		var data interface{}

		code = err.SUCCESS           // 设置 code 为成功
		token := cont.Query("token") // 获取查询参数中的 token

		// 如果 token 为空，设置 code 为无效参数
		if token == "" {
			code = err.INVALID_PARAMS
		} else {
			// 解析 token
			claims, e := util.ParseToken(token)

			// 如果解析失败，设置 code 为解析失败
			if e != nil {
				code = err.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				// 如果 token 过期，设置 code 为 token 过期
				code = err.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		// 如果 code 不为成功，返回 json 响应
		if code != err.SUCCESS {
			cont.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  err.GetMsg(code),
				"data": data,
			})
			// 终止处理
			cont.Abort()
			return
		}

		// 继续处理请求
		cont.Next()
	}
}
