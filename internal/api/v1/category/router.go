package category

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册分类相关路由
func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/categories", ctrl.listCategories)
}
