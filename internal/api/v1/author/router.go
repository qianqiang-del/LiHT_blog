package author

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册作者相关路由
func (ctrl *Controller) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/author", ctrl.getAuthorDetail)
	rg.GET("/about", ctrl.getAbout)
}
