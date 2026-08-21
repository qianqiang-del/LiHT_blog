package category

import (
	"blog/internal/service"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// Controller 分类控制器
type Controller struct {
	svc service.CategoryService
}

// NewController 创建分类控制器
func NewController(svc service.CategoryService) *Controller {
	return &Controller{svc: svc}
}

// listCategories GET /api/v1/categories
func (ctrl *Controller) listCategories(c *gin.Context) {
	result, err := ctrl.svc.GetCategoryDetail()
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, result)
}
