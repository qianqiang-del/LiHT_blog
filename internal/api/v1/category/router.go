package category

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册分类相关路由
func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/categories", ctrl.listCategories)
	r.GET("/categories/:id/articles", ctrl.listCategoryArticles)
}

// RegisterAdminRoutes 注册后台分类管理路由
func (ctrl *Controller) RegisterAdminRoutes(r *gin.RouterGroup) {
	r.GET("/categories", ctrl.listCategories)
	r.POST("/categories", ctrl.adminCreateCategory)
	r.DELETE("/categories/:id", ctrl.adminDeleteCategory)
}
