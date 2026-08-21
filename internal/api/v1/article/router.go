package article

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册文章相关路由
func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/articles", ctrl.listArticles)
	r.GET("/hot", ctrl.listHotArticles)
	r.GET("/articles/:id", ctrl.getArticleDetail)
}
