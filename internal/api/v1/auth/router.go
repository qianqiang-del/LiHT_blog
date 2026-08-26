package auth

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册认证相关路由
func (ctrl *Controller) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/auth/send-code", ctrl.sendCode)
	rg.POST("/auth/register", ctrl.register)
	rg.POST("/auth/login", ctrl.login)
	rg.POST("/auth/logout", ctrl.logout)
	rg.GET("/auth/captcha", ctrl.getCaptcha)
}
