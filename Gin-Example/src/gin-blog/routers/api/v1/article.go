package v1

import (
	"Gin-Example/src/gin-blog/models"
	"Gin-Example/src/gin-blog/pkg/err"
	"Gin-Example/src/gin-blog/pkg/setting"
	"Gin-Example/src/gin-blog/pkg/util"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

// GetArticle 获取单个文章
func GetArticle(cont *gin.Context) {
	id := com.StrTo(cont.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := err.INVALID_PARAMS
	var data interface{}
	if !valid.HasErrors() {
		if models.ExistTagByID(id) {
			data = models.GetArticle(id)
			code = err.SUCCESS
		} else {
			code = err.ERROR_NOT_EXIST_ARTICLE
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

// GetArticles 获取多个文章
func GetArticles(cont *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	valid := validation.Validation{}

	var state int = -1
	if arg := cont.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state

		valid.Range(state, 0, 1, "state").Message("状态只允许为0 / 1")
	}

	var tagID int = -1
	if arg := cont.Query("tag_id"); arg != "" {
		tagID = com.StrTo(arg).MustInt()
		maps["tag_id"] = tagID

		valid.Min(tagID, 1, "tag_id").Message("标签ID必须大于0")
	}

	code := err.INVALID_PARAMS
	if !valid.HasErrors() {
		code = err.SUCCESS

		data["lists"] = models.GetArticles(util.GetPage(cont), setting.PageSize, maps)
		data["total"] = models.GetArticleTotal(maps)
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

// AddArticle 新增文章
func AddArticle(cont *gin.Context) {
	tagID := com.StrTo(cont.Query("tag_id")).MustInt()
	title := cont.Query("title")
	desc := cont.Query("desc")
	content := cont.Query("content")
	createdBy := cont.Query("created_by")
	state := com.StrTo(cont.DefaultQuery("state", "0")).MustInt()

	valid := validation.Validation{}
	valid.Min(tagID, 1, "tag_id").Message("标签ID必须大于0")
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(desc, "desc").Message("简述不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(createdBy, "created_by").Message("创建人信息不能为空")
	valid.Range(state, 0, 1, "state").Message("状态只允许为0 / 1")

	code := err.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistTagByID(tagID) {
			data := make(map[string]interface{})
			data["tag_id"] = tagID
			data["title"] = title
			data["desc"] = desc
			data["content"] = content
			data["created_by"] = createdBy
			data["state"] = state

			models.AddArticle(data)
			code = err.SUCCESS
		} else {
			code = err.ERROR_NOT_EXIST_TAG
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

// UpdateArticle 修改文章
func UpdateArticle(cont *gin.Context) {
	valid := validation.Validation{}

	id := com.StrTo(cont.Param("id")).MustInt()
	tagId := com.StrTo(cont.Query("tag_id")).MustInt()
	title := cont.Query("title")
	desc := cont.Query("desc")
	content := cont.Query("content")
	modifiedBy := cont.Query("modified_by")

	var state int = -1
	if arg := cont.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state")
	}

	valid.Min(id, 1, "id").Message("ID必须大于0")
	valid.MaxSize(title, 128, "title").Message("标题最长为128个字符")
	valid.MaxSize(desc, 255, "desc").Message("简述最多为255个字符")
	valid.MaxSize(content, 65535, "content").Message("内容最长为65535")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 128, "modified_bt").Message("修改人最长为128字符")

	code := err.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			data := make(map[string]interface{})
			if tagId > 0 {
				data["tag_id"] = tagId
			}
			if title != "" {
				data["title"] = title
			}
			if desc != "" {
				data["desc"] = desc
			}
			if content != "" {
				data["content"] = content
			}
			data["modified_by"] = modifiedBy
			models.EditArticle(id, data)
			code = err.SUCCESS
		} else {
			code = err.ERROR_NOT_EXIST_TAG
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

// DeleteArticle 删除文章
func DeleteArticle(cont *gin.Context) {
	id := com.StrTo(cont.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := err.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			models.DeleteArticle(id)
			code = err.SUCCESS
		} else {
			code = err.ERROR_NOT_EXIST_ARTICLE
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
