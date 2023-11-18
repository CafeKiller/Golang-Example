package main

import (
	"Gin-Example/src/gin-blog/pkg/setting"
	"Gin-Example/src/gin-blog/routers"
	"fmt"

	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {

	router := routers.InitRouter()

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	server.ListenAndServe()

}

func GetAppPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	return path[:index]
}
