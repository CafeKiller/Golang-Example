package routers

import (
	"Gin-Example/src/gin-blog/pkg/setting"
	v1 "Gin-Example/src/gin-blog/routers/api/v1"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {

	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	r.GET("/hello", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "Hello GO! Hello Gin! Hi CoffeeKiller",
		})
	})

	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/tags", v1.GetTags)
		apiV1.POST("/tags", v1.AddTag)
		apiV1.PUT("/tags/:id", v1.EditTag)
		apiV1.DELETE("/tags/:id", v1.DeleteTag)
	}

	return r
}
