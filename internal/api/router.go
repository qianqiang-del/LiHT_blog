package api

import (
	"blog/internal/api/v1/article"
	"blog/internal/api/v1/auth"
	"blog/internal/api/v1/author"
	"blog/internal/api/v1/category"
	"blog/internal/api/v1/comment"
	"blog/internal/api/v1/tag"
	"blog/internal/api/v1/upload"
	"blog/internal/api/v1/user"
	"blog/internal/middleware"
	"blog/internal/repository"
	"blog/internal/service"

	"github.com/gin-gonic/gin"
)

// Router 路由
type Router struct {
	articleCtrl  *article.Controller
	authorCtrl   *author.Controller
	categoryCtrl *category.Controller
	tagCtrl      *tag.Controller
	authCtrl     *auth.Controller
	commentCtrl  *comment.Controller
	uploadCtrl   *upload.Controller
	userCtrl     *user.Controller
}

// NewRouter 创建路由
func NewRouter(
	articleSvc service.ArticleService,
	authorSvc service.AuthorService,
	categorySvc service.CategoryService,
	tagSvc service.TagService,
	authSvc service.AuthService,
	commentSvc service.CommentService,
	uploadSvc service.UploadService,
	userSvc service.UserService,
) *Router {
	return &Router{
		articleCtrl:  article.NewController(articleSvc),
		authorCtrl:   author.NewController(authorSvc),
		categoryCtrl: category.NewController(categorySvc),
		tagCtrl:      tag.NewController(tagSvc),
		authCtrl:     auth.NewController(authSvc),
		commentCtrl:  comment.NewController(commentSvc),
		uploadCtrl:   upload.NewController(uploadSvc),
		userCtrl:     user.NewController(userSvc),
	}
}

// Setup 设置路由
func (r *Router) Setup(engine *gin.Engine, authRepo repository.AuthRepository) {
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

	// 前台 API 路由组
	apiGroup := engine.Group("/api/v1")
	{
		r.articleCtrl.RegisterRoutes(apiGroup, authRepo)
		r.authorCtrl.RegisterRoutes(apiGroup)
		r.categoryCtrl.RegisterRoutes(apiGroup)
		r.tagCtrl.RegisterRoutes(apiGroup)
		r.authCtrl.RegisterRoutes(apiGroup, authRepo)
		r.commentCtrl.RegisterRoutes(apiGroup, authRepo)
	}

	// 后台管理路由组
	adminGroup := engine.Group("/api/v1/admin")
	{
		r.articleCtrl.RegisterAdminRoutes(adminGroup)
		r.categoryCtrl.RegisterAdminRoutes(adminGroup)
		r.tagCtrl.RegisterAdminRoutes(adminGroup)
		r.commentCtrl.RegisterAdminRoutes(adminGroup)
		r.uploadCtrl.RegisterRoutes(adminGroup)
		r.userCtrl.RegisterAdminRoutes(adminGroup)
	}
}

// Close 关闭所有路由连接
func (r *Router) Close() error {
	return nil
}
