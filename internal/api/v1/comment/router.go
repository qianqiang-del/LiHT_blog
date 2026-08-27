package comment

import (
	"blog/internal/middleware"
	"blog/internal/repository"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册评论相关路由
func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup, authRepo repository.AuthRepository) {
	// 获取评论列表（可选登录，用于判断是否点赞）
	r.GET("/articles/:id/comments", middleware.OptionalAuth(authRepo), ctrl.listComments)

	// 获取回复列表（可选登录）
	r.GET("/comments/:id/replies", middleware.OptionalAuth(authRepo), ctrl.listReplies)

	// 需要登录的路由
	authGroup := r.Group("")
	authGroup.Use(middleware.Auth(authRepo))
	// 创建评论
	authGroup.POST("/articles/:id/comments", ctrl.createComment)
	// 点赞/取消点赞评论
	authGroup.POST("/comments/:id/like", ctrl.likeComment)
}
