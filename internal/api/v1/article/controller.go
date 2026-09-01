package article

import (
	"blog/internal/middleware"
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

// searchArticles GET /api/v1/articles/search
func (ctrl *Controller) searchArticles(c *gin.Context) {
	var req request.ArticleSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "搜索参数错误")
		return
	}

	result, err := ctrl.svc.SearchArticles(req)
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

	// OptionalAuth：未登录 userID 为 0，登录了有值
	var userID *uint
	if uid := middleware.GetUserID(c); uid > 0 {
		userID = &uid
	}

	result, err := ctrl.svc.GetArticleDetail(uint(id), userID)
	if err != nil {
		if bizErr, ok := errors.Is(err, errors.CodeResourceNotFound); ok {
			response.NotFound(c, bizErr.Message)
			return
		}
		response.BizError(c, err)
		return
	}

	response.Success(c, result)
}

// likeArticle POST /api/v1/articles/:id/like
func (ctrl *Controller) likeArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的文章 ID")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}

	result, err := ctrl.svc.LikeArticle(uint(id), userID)
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, result)
}

// adminListArticles GET /api/v1/admin/articles
func (ctrl *Controller) adminListArticles(c *gin.Context) {
	var req request.AdminArticleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "分页参数错误")
		return
	}
	result, err := ctrl.svc.AdminListArticles(req)
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, result)
}

// adminGetArticleDetail GET /api/v1/admin/articles/:id
func (ctrl *Controller) adminGetArticleDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的文章 ID")
		return
	}

	result, err := ctrl.svc.GetArticleDetail(uint(id), nil)
	if err != nil {
		if bizErr, ok := errors.Is(err, errors.CodeResourceNotFound); ok {
			response.NotFound(c, bizErr.Message)
			return
		}
		response.BizError(c, err)
		return
	}

	response.Success(c, result)
}

// adminUpdateArticleStatus PATCH /api/v1/admin/articles/:id/status
func (ctrl *Controller) adminUpdateArticleStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的文章 ID")
		return
	}

	var req request.AdminUpdateArticleStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "状态值无效，只能为 正常 或 下架")
		return
	}

	// 转换为数据库值：正常=1, 下架=0
	var status int8
	if req.Status == "正常" {
		status = 1
	} else {
		status = 0
	}

	if err := ctrl.svc.AdminUpdateArticleStatus(uint(id), status); err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, nil)
}

// adminDeleteArticle DELETE /api/v1/admin/articles/:id
func (ctrl *Controller) adminDeleteArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的文章 ID")
		return
	}

	if err := ctrl.svc.AdminDeleteArticle(uint(id)); err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, nil)
}

// adminCreateArticle POST /api/v1/admin/articles
func (ctrl *Controller) adminCreateArticle(c *gin.Context) {
	var req request.AdminCreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	if err := ctrl.svc.AdminCreateArticle(req); err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, "发布成功")
}

// adminUpdateArticle PUT /api/v1/admin/articles/:id
func (ctrl *Controller) adminUpdateArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的文章 ID")
		return
	}

	var req request.AdminUpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	if err := ctrl.svc.AdminUpdateArticle(uint(id), req); err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, "更新成功")
}

// adminListArticleOptions GET /api/v1/admin/articles/options
func (ctrl *Controller) adminListArticleOptions(c *gin.Context) {
	result, err := ctrl.svc.AdminListArticleOptions()
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, result)
}
