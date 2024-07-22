package util

import (
	"Gin-Example/src/gin-blog/pkg/setting"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// GetPage 获取分页页码
func GetPage(c *gin.Context) int {

	// 初始化result为0
	result := 0

	// 从查询参数中获取page，并将其转换为int类型
	page, _ := com.StrTo(c.Query("page")).Int()

	// 如果page大于0，则将page减1，再乘以setting.PageSize，得到result
	if page > 0 {
		result = (page - 1) * setting.PageSize
	}

	return result
}
