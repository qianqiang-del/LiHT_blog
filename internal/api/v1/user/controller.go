package user

import (
	"blog/internal/model/dto/request"
	"blog/internal/service"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// Controller 用户控制器
type Controller struct {
	svc service.UserService
}

// NewController 创建用户控制器
func NewController(svc service.UserService) *Controller {
	return &Controller{svc: svc}
}

// adminListUsers GET /api/v1/admin/users
func (ctrl *Controller) adminListUsers(c *gin.Context) {
	var req request.AdminUserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "分页参数错误")
		return
	}

	result, err := ctrl.svc.AdminListUsers(req.Page, req.Size)
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, result)
}
