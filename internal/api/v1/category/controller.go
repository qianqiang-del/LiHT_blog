package category

import (
	"blog/internal/model/dto/request"
	"blog/internal/service"
	"blog/pkg/errors"
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

// adminCreateCategory POST /api/v1/admin/categories
func (ctrl *Controller) adminCreateCategory(c *gin.Context) {
	var req request.AdminCreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "分类名称不能为空且不超过20字")
		return
	}

	if err := ctrl.svc.AdminCreateCategory(req); err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, "添加成功")
}

// adminDeleteCategory DELETE /api/v1/admin/categories/:id
func (ctrl *Controller) adminDeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的分类 ID")
		return
	}

	if err := ctrl.svc.AdminDeleteCategory(uint(id)); err != nil {
		if bizErr, ok := errors.Is(err, errors.CodeBadRequest); ok {
			response.BadRequest(c, bizErr.Message)
			return
		}
		response.BizError(c, err)
		return
	}

	response.Success(c, "删除成功")
}
