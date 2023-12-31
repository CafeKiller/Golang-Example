package routers

import (
	"Gin-Example/src/gin-blog/middleware/jwt"
	"Gin-Example/src/gin-blog/pkg/setting"
	"Gin-Example/src/gin-blog/routers/api"
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

	r.GET("/auth", api.GetAuth)

	apiV1 := r.Group("/api/v1")
	apiV1.Use(jwt.JWT())

	{
		apiV1.GET("/tags", v1.GetTags)
		apiV1.POST("/tags", v1.AddTag)
		apiV1.PUT("/tags/:id", v1.EditTag)
		apiV1.DELETE("/tags/:id", v1.DeleteTag)

		apiV1.GET("/articles", v1.GetArticles)
		apiV1.GET("/articles/:id", v1.GetArticle)
		apiV1.POST("/articles", v1.AddArticle)
		apiV1.PUT("/articles/:id", v1.EditTag)
		apiV1.DELETE("/articles/:id", v1.DeleteArticle)
	}

	return r
}
