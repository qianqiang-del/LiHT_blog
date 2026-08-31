package user

import (
	"blog/internal/model/dto/request"
	"blog/internal/service"
	"blog/pkg/response"
	"strconv"

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

// adminUpdateUserStatus PUT /api/v1/admin/users/:id/status
func (ctrl *Controller) adminUpdateUserStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户 ID")
		return
	}

	var req request.AdminUpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "状态参数错误")
		return
	}

	if err := ctrl.svc.AdminUpdateUserStatus(uint(id), req); err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, nil)
}

// adminDeleteUser DELETE /api/v1/admin/users/:id
func (ctrl *Controller) adminDeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户 ID")
		return
	}

	if err := ctrl.svc.AdminDeleteUser(uint(id)); err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, nil)
}
