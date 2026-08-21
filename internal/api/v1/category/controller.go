package category

import (
	"blog/internal/model/dto/request"
	"blog/internal/service"
	"blog/pkg/response"
	"strconv"

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

// listCategoryArticles GET /api/v1/categories/:id/articles
func (ctrl *Controller) listCategoryArticles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的分类 ID")
		return
	}

	var req request.ArticleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "分页参数错误")
		return
	}

	result, err := ctrl.svc.ListCategoryArticles(uint(id), req)
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, result)
}
