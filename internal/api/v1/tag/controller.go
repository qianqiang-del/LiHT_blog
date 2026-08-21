package tag

import (
	"blog/internal/service"
	"blog/pkg/response"

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
func (ctrl *Controller) listTags(c *gin.Context) {
	result, err := ctrl.svc.ListTags()
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, result)
}
