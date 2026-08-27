package auth

import (
	"blog/internal/middleware"
	"blog/internal/repository"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册认证相关路由
func (ctrl *Controller) RegisterRoutes(rg *gin.RouterGroup, authRepo repository.AuthRepository) {
	rg.POST("/auth/send-code", ctrl.sendCode)
	rg.POST("/auth/register", ctrl.register)
	rg.POST("/auth/login", ctrl.login)
	rg.POST("/auth/logout", ctrl.logout)
	rg.GET("/auth/captcha", ctrl.getCaptcha)

	// 需要登录的接口
	authGroup := rg.Group("/auth")
	authGroup.Use(middleware.Auth(authRepo))
	authGroup.GET("/me", ctrl.me)
}
