package v1

import (
	"Gin-Example/src/gin-blog/models"
	"Gin-Example/src/gin-blog/pkg/err"
	"Gin-Example/src/gin-blog/pkg/setting"
	"Gin-Example/src/gin-blog/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

// GetTags 获取多个文章标签
func GetTags(cont *gin.Context) {
	// 获取URL参数
	name := cont.Query("name")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}

	var state int = -1
	if arg := cont.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	code := err.SUCCESS

	data["lists"] = models.GetTags(util.GetPage(cont), setting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)

	cont.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  err.GetMsg(code),
		"data": data,
	})
}

// AddTag 添加一个文章标签
func AddTag(cont *gin.Context) {

}

// EditTag 修改文章标签
func EditTag(cont *gin.Context) {

}

// DeleteTag 删除文章标签
func DeleteTag(cont *gin.Context) {

}
