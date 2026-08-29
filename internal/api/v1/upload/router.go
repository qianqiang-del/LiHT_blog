package upload

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册上传相关路由
func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/upload", ctrl.upload)
}
