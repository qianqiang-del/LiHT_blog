package article

import (
	"blog/internal/model/dto/request"
	"blog/internal/service"
	"blog/pkg/errors"
	"blog/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Controller 文章控制器
type Controller struct {
	svc service.ArticleService
}

// NewController 创建文章控制器
func NewController(svc service.ArticleService) *Controller {
	return &Controller{svc: svc}
}

// listArticles GET /api/v1/articles
func (ctrl *Controller) listArticles(c *gin.Context) {
	var req request.ArticleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "分页参数错误")
		return
	}
	result, err := ctrl.svc.ListArticles(req)
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, result)
}

// listHotArticles GET /api/v1/articles/hot
func (ctrl *Controller) listHotArticles(c *gin.Context) {
	var req request.ArticleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "分页参数错误")
		return
	}
	result, err := ctrl.svc.ListHotArticles(req)
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, result)
}

// getArticleDetail GET /api/v1/articles/:id
func (ctrl *Controller) getArticleDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的文章 ID")
		return
	}

	result, err := ctrl.svc.GetArticleDetail(uint(id))
	if err != nil {
		if bizErr, ok := err.(*errors.BizError); ok && bizErr.Code == errors.CodeResourceNotFound {
			response.NotFound(c, bizErr.Message)
			return
		}
		response.BizError(c, err)
		return
	}

	response.Success(c, result)
}
