package dashboard

import "github.com/gin-gonic/gin"

// RegisterAdminRoutes 注册后台概览路由
func (ctrl *Controller) RegisterAdminRoutes(r *gin.RouterGroup) {
	r.GET("/dashboard", ctrl.adminGetDashboard)
}
