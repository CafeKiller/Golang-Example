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

	// 创建一个gin引擎
	r := gin.New()

	// 使用gin的日志记录功能
	r.Use(gin.Logger())

	// 使用gin的错误恢复功能
	r.Use(gin.Recovery())

	// 设置运行模式
	gin.SetMode(setting.RunMode)

	// 测试接口
	r.GET("/hello", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "Hello GO! Hello Gin! Hi CoffeeKiller",
		})
	})

	// 认证接口
	r.GET("/auth", api.GetAuth)

	// 创建一个以/api/v1为前缀的路由分组
	apiV1 := r.Group("/api/v1")

	// 使用jwt中间件
	apiV1.Use(jwt.JWT())

	// 定义一个匿名路由，用于处理未匹配的路由
	{
		// 获取标签
		apiV1.GET("/tags", v1.GetTags)
		// 添加标签
		apiV1.POST("/tags", v1.AddTag)
		// 编辑标签
		apiV1.PUT("/tags/:id", v1.EditTag)
		// 删除标签
		apiV1.DELETE("/tags/:id", v1.DeleteTag)

		// 获取文章
		apiV1.GET("/articles", v1.GetArticles)
		// 获取文章详情
		apiV1.GET("/articles/:id", v1.GetArticle)
		// 添加文章
		apiV1.POST("/articles", v1.AddArticle)
		// 编辑文章
		apiV1.PUT("/articles/:id", v1.EditTag)
		// 删除文章
		apiV1.DELETE("/articles/:id", v1.DeleteArticle)
	}

	return r
}
