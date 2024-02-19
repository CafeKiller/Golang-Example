package main

import (
	"github.com/labstack/echo"
	"html/template"
	"io"
)

// Template HTML 渲染模版对象
type Template struct { // [review] 此处为一个接口(Renderer接口)的实现
}

// Render Template 对象的接口, 将数据写入 writer, 并嵌入HTML模块
func (t *Template) Render(writer io.Writer, name string, data interface{}, content echo.Context) error {

	// 通过 name 读取全局模版集合中指定模版
	if t, ok := templates[name]; ok {
		return t.ExecuteTemplate(writer, "layout.html", data)
	}
	content.Echo().Logger.Debugf("Template[%s] Not Found", name)
	return templates["error"].ExecuteTemplate(writer, "layout.html", "Internal Server Error")
}

// loadTemplates 读取HTML模版
func loadTemplates() {
	// 基础模版
	var baseTemplate = "webserver/templates/layout.html"

	// 将所有html模版逐个放入 全局对象: templates 集合当中
	templates = make(map[string]*template.Template)
	/*
		template.ParseFiles 表示通过文件的方式解析模版
		template.ParseFiles 可以传入多个string参数,
		其返回的模版是第一个传入模版, 后续模版则是用于解析, 解析后的结果传递到第一个模版
	*/
	templates["index"] = template.Must(
		template.ParseFiles(baseTemplate, "webserver/templates/index.html"))
	templates["error"] = template.Must(
		template.ParseFiles(baseTemplate, "webserver/templates/error.html"))
	templates["user"] = template.Must(
		template.ParseFiles(baseTemplate, "webserver/templates/user.html"))
	templates["login"] = template.Must(
		template.ParseFiles(baseTemplate, "webserver/templates/login.html"))
	templates["admin"] = template.Must(
		template.ParseFiles(baseTemplate, "webserver/templates/admin.html"))
	templates["admin_users"] = template.Must(
		template.ParseFiles(baseTemplate, "webserver/templates/admin_users.html"))
}
