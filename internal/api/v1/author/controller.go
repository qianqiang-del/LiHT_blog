package author

import (
	"blog/internal/model/dto/request"
	"blog/internal/service"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// Controller 作者控制器
type Controller struct {
	svc service.AuthorService
}

// NewController 创建作者控制器
func NewController(svc service.AuthorService) *Controller {
	return &Controller{svc: svc}
}

// getAuthorDetail GET /api/v1/author
func (ctrl *Controller) getAuthorDetail(c *gin.Context) {
	result, err := ctrl.svc.GetAuthorDetail()
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, result)
}

// getAbout GET /api/v1/about
func (ctrl *Controller) getAbout(c *gin.Context) {
	result, err := ctrl.svc.GetAbout()
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, result)
}

// adminGetAuthor GET /api/v1/admin/author
func (ctrl *Controller) adminGetAuthor(c *gin.Context) {
	result, err := ctrl.svc.GetAbout()
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, result)
}

// adminUpdateAuthor PUT /api/v1/admin/author
func (ctrl *Controller) adminUpdateAuthor(c *gin.Context) {
	var req request.AdminUpdateAuthorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	if err := ctrl.svc.AdminUpdateAuthor(req); err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, nil)
}
