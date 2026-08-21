package tag

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册标签相关路由
func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/tags", ctrl.listTags)
}
