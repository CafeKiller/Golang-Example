package main

import (
	"github.com/labstack/echo"
	"html/template"
	"io"
)

// Template HTML渲染模版
type Template struct {
}

// Render 将数据写入writer, 并嵌入HTML模块
func (t *Template) Render(writer io.Writer, name string, data interface{}, content echo.Context) error {

	if t, ok := templates[name]; ok {
		return t.ExecuteTemplate(writer, "layout.html", data)
	}

	content.Echo().Logger.Debugf("Template[%s] Not Found", name)

	return templates["error"].ExecuteTemplate(writer, "layout.html", "Internal Server Error")
}

// loadTemplates 读取HTML模版
func loadTemplates() {
	var baseTemplate = "webserver/templates/layout.html"
	templates = make(map[string]*template.Template)
	// 将所有html模版放入templates集合中
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
