package user

import "github.com/gin-gonic/gin"

// RegisterAdminRoutes 注册后台用户管理路由
func (ctrl *Controller) RegisterAdminRoutes(r *gin.RouterGroup) {
	r.GET("/users", ctrl.adminListUsers)
}
