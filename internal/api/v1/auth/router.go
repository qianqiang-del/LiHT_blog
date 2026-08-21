package auth

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册认证相关路由
func (ctrl *Controller) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/auth/send-code", ctrl.sendCode)
	rg.POST("/auth/register", ctrl.register)
}
