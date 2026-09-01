package author

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册作者相关路由
func (ctrl *Controller) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/author", ctrl.getAuthorDetail)
	rg.GET("/about", ctrl.getAbout)
}

// RegisterPublicRoutes 注册无需鉴权的后台路由
func (ctrl *Controller) RegisterPublicRoutes(rg *gin.RouterGroup) {
	rg.POST("/login", ctrl.adminLogin)
	rg.GET("/captcha", ctrl.adminGetCaptcha)
}

// RegisterAdminRoutes 注册需要鉴权的后台路由
func (ctrl *Controller) RegisterAdminRoutes(rg *gin.RouterGroup) {
	rg.POST("/logout", ctrl.adminLogout)
	rg.GET("/author", ctrl.adminGetAuthor)
	rg.PUT("/author", ctrl.adminUpdateAuthor)
	rg.PUT("/author/account", ctrl.adminUpdateAccount)
	rg.PUT("/author/password", ctrl.adminUpdatePassword)
}
