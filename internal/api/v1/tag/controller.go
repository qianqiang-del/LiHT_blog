package tag

import (
	"blog/internal/model/dto/request"
	"blog/internal/service"
	"blog/pkg/response"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Controller 标签控制器
type Controller struct {
	svc service.TagService
}

// NewController 创建标签控制器
func NewController(svc service.TagService) *Controller {
	return &Controller{svc: svc}
}

// listTags GET /api/v1/tags

// listTags GET /api/v1/tags
func (ctrl *Controller) listTags(c *gin.Context) {
	result, err := ctrl.svc.ListTags()
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, result)
}

// listTagArticles GET /api/v1/tags/:id/articles
func (ctrl *Controller) listTagArticles(c *gin.Context) {
	// 解析逗号分隔的 ID：1,2,3
	idStr := c.Param("id")
	parts := strings.Split(idStr, ",")
	var tagIDs []uint
	for _, p := range parts {
		id, err := strconv.ParseUint(strings.TrimSpace(p), 10, 64)
		if err != nil {
			response.BadRequest(c, "无效的标签 ID")
			return
		}
		tagIDs = append(tagIDs, uint(id))
	}

	var req request.ArticleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "分页参数错误")
		return
	}

	result, err := ctrl.svc.ListTagArticles(tagIDs, req)
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, result)
}

// adminCreateTag POST /api/v1/admin/tags
func (ctrl *Controller) adminCreateTag(c *gin.Context) {
	var req request.AdminCreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "标签名称不能为空且不超过20字")
		return
	}

	if err := ctrl.svc.AdminCreateTag(req.Name); err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, "添加成功")
}
