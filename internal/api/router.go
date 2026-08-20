package api

import (
	"blog/internal/api/v1/article"
	"blog/internal/api/v1/author"
	"blog/internal/middleware"
	"blog/internal/service"

	"github.com/gin-gonic/gin"
)

// Router 路由
type Router struct {
	articleCtrl *article.Controller
	authorCtrl  *author.Controller
}

// NewRouter 创建路由
func NewRouter(
	articleSvc service.ArticleService,
	authorSvc service.AuthorService,
) *Router {
	return &Router{
		articleCtrl: article.NewController(articleSvc),
		authorCtrl:  author.NewController(authorSvc),
	}
}

// Setup 设置路由
func (r *Router) Setup(engine *gin.Engine) {
	// 全局中间件
	engine.Use(middleware.Recovery())
	engine.Use(middleware.Logger())
	engine.Use(middleware.CORS())

	// 健康检查
	engine.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "LiHT API is running",
		})
	})

	// API 路由组
	apiGroup := engine.Group("/api/v1")
	{
		r.articleCtrl.RegisterRoutes(apiGroup)
		r.authorCtrl.RegisterRoutes(apiGroup)
	}
}

// Close 关闭所有路由连接
func (r *Router) Close() error {
	return nil
}
