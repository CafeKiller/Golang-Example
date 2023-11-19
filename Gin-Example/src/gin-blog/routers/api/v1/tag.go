package v1

import (
	"github.com/astaxie/beego/logs"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"Gin-Example/src/gin-blog/models"
	"Gin-Example/src/gin-blog/pkg/err"
	"Gin-Example/src/gin-blog/pkg/setting"
	"Gin-Example/src/gin-blog/pkg/util"
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
	name := cont.Query("name")
	state := com.StrTo(cont.DefaultQuery("state", "0")).MustInt()
	createdBy := cont.Query("created_by")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许为0/1")

	code := err.INVALID_PARAMS

	if !valid.HasErrors() {
		if !models.ExistTagByName(name) {
			code = err.SUCCESS
			models.AddTag(name, state, createdBy)
		} else {
			code = err.ERROR_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			logs.Info(err.Key, err.Message)
		}
	}

	cont.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  err.GetMsg(code),
		"data": make(map[string]string),
	})

}

// EditTag 修改文章标签
func EditTag(cont *gin.Context) {

}

// DeleteTag 删除文章标签
func DeleteTag(cont *gin.Context) {

}
