package auth

import (
	"blog/internal/model/dto/request"
	"blog/internal/service"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// Controller 认证控制器
type Controller struct {
	svc service.AuthService
}

// NewController 创建认证控制器
func NewController(svc service.AuthService) *Controller {
	return &Controller{svc: svc}
}

// sendCode POST /api/v1/auth/send-code
func (ctrl *Controller) sendCode(c *gin.Context) {
	var req request.SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "邮箱格式错误")
		return
	}

	if err := ctrl.svc.SendCode(req); err != nil {
		response.BizError(c, err)
		return
	}

	response.SuccessWithMessage(c, "验证码已发送", nil)
}

// register POST /api/v1/auth/register
func (ctrl *Controller) register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.svc.Register(req); err != nil {
		response.BizError(c, err)
		return
	}

	response.SuccessWithMessage(c, "注册成功", nil)
}
