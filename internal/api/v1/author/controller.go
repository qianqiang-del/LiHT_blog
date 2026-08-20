package author

import (
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
