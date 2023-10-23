package main

import (
	"github.com/labstack/echo"
	"io"
)

// Template HTML渲染模版
type Template struct {
}

func (t *Template) Render(writer io.Writer, name string, data interface{}, content echo.Context) error {

	if t, ok := templates[name]; ok {
		return t.ExecuteTemplate(writer, "layout.html", data)
	}

}
