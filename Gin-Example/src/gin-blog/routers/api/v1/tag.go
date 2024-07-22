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

	// 初始化maps和data
	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	// 如果name不为空，则将其添加到maps中
	if name != "" {
		maps["name"] = name
	}

	// 初始化 state
	var state int = -1
	// 如果 arg 不为空，则将其转换为int类型，并将其添加到maps中
	if arg := cont.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	// 初始化 code
	code := err.SUCCESS

	// 获取文章标签列表和总数，并将它们添加到 data 中
	data["lists"] = models.GetTags(util.GetPage(cont), setting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)

	// 返回 JSON 数据
	cont.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  err.GetMsg(code),
		"data": data,
	})
}

// AddTag @Summary 新增文章标签
// @Produce  json
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param created_by query int false "CreatedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags [post]
func AddTag(cont *gin.Context) {
	// 获取参数
	name := cont.Query("name")
	state := com.StrTo(cont.DefaultQuery("state", "0")).MustInt()
	createdBy := cont.Query("created_by")

	// 验证参数
	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许为0/1")

	// 初始化返回码
	code := err.INVALID_PARAMS

	// 验证参数是否有效
	if !valid.HasErrors() {
		// 检查标签是否存在
		if !models.ExistTagByName(name) {
			// 标签不存在，新增标签
			code = err.SUCCESS
			models.AddTag(name, state, createdBy)
		} else {
			// 标签已存在
			code = err.ERROR_EXIST_TAG
		}
	} else {
		// 验证参数出错
		for _, err := range valid.Errors {
			logs.Info(err.Key, err.Message)
		}
	}

	// 返回结果
	cont.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  err.GetMsg(code),
		"data": make(map[string]string),
	})

}

// EditTag @Summary 修改文章标签
// @Produce  json
// @Param id param int true "ID"
// @Param name query string true "ID"
// @Param state query int false "State"
// @Param modified_by query string true "ModifiedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags/{id} [put]
func EditTag(cont *gin.Context) {
	// 从路径参数中获取ID
	id := com.StrTo(cont.Param("id")).MustInt()
	// 从查询参数中获取name
	name := cont.Query("name")
	// 从查询参数中获取modified_by
	modifiedBy := cont.Query("modified_by")

	// 验证参数
	valid := validation.Validation{}

	// 从查询参数中获取state，如果没有则默认为-1
	var state int = -1
	if arg := cont.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许为0 / 1")
	}

	// 验证必填项
	valid.Required(id, "id").Message("ID不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(name, 127, "name").Message("名称最长为127个字符")
	valid.MaxSize(modifiedBy, 255, "modified_by").Message("修改人最长255字符")

	// 如果验证通过，则code为200，否则为400
	code := err.INVALID_PARAMS

	// 如果验证有误，则打印错误信息
	print("id=====", id)

	if !valid.HasErrors() {
		code = err.SUCCESS
		// 如果ID存在，则编辑标签
		if models.ExistTagByID(id) {
			data := make(map[string]interface{})
			data["modified_by"] = modifiedBy
			if name != "" {
				data["name"] = name
			}
			if state != -1 {
				data["state"] = state
			}

			models.EditTag(id, data)
		} else {
			// 如果ID不存在，则打印错误信息
			for _, err := range valid.Errors {
				logs.Info(err.Key, err.Message)
			}
		}
	}

	// 返回结果
	cont.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  err.GetMsg(code),
		"data": make(map[string]string),
	})

}

// DeleteTag 删除文章标签
func DeleteTag(cont *gin.Context) {
	// 从URL中获取id参数
	id := com.StrTo(cont.Param("id")).MustInt()

	// 验证id参数
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	// 定义错误码
	code := err.INVALID_PARAMS

	// 验证参数是否有效
	if !valid.HasErrors() {
		code = err.SUCCESS
		// 检查标签是否存在
		if models.ExistTagByID(id) {
			models.DeleteTag(id) // 删除标签
		} else {
			code = err.ERROR_EXIST_TAG
		}
	} else {
		// 记录验证错误信息
		for _, err := range valid.Errors {
			logs.Info(err.Key, err.Message)
		}
	}

	// 返回结果
	cont.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  err.GetMsg(code),
		"data": make(map[string]string),
	})
}
