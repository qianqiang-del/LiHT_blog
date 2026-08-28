package article

import (
	"blog/internal/middleware"
	"blog/internal/repository"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册文章相关路由
func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup, authRepo repository.AuthRepository) {
	r.GET("/articles", ctrl.listArticles)
	r.GET("/articles/search", ctrl.searchArticles)
	r.GET("/hot", ctrl.listHotArticles)
	r.GET("/articles/:id", middleware.OptionalAuth(authRepo), ctrl.getArticleDetail)

	// 需要登录的路由
	authGroup := r.Group("/articles")
	authGroup.Use(middleware.Auth(authRepo))
	authGroup.POST("/:id/like", ctrl.likeArticle)
}

// RegisterAdminRoutes 注册后台文章管理路由
func (ctrl *Controller) RegisterAdminRoutes(r *gin.RouterGroup) {
	r.GET("/articles", ctrl.adminListArticles)
	r.GET("/articles/:id", ctrl.adminGetArticleDetail)
	r.PUT("/articles/:id/status", ctrl.adminUpdateArticleStatus)
	r.DELETE("/articles/:id", ctrl.adminDeleteArticle)
}
