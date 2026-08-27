package auth

import (
	"blog/internal/middleware"
	"blog/internal/model/dto/request"
	"blog/internal/service"
	"blog/pkg/response"
	"strings"

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

// login POST /api/v1/auth/login
func (ctrl *Controller) login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	result, err := ctrl.svc.Login(req)
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, result)
}

// getCaptcha GET /api/v1/auth/captcha
func (ctrl *Controller) getCaptcha(c *gin.Context) {
	result, err := ctrl.svc.GenerateCaptcha()
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, result)
}

// me GET /api/v1/auth/me
func (ctrl *Controller) me(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}

	user, err := ctrl.svc.GetUserInfo(userID)
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, user)
}

// logout POST /api/v1/auth/logout
func (ctrl *Controller) logout(c *gin.Context) {
	// 从 Header 取 token
	authHeader := c.GetHeader("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		response.BadRequest(c, "令牌格式错误")
		return
	}

	if err := ctrl.svc.Logout(parts[1]); err != nil {
		response.BizError(c, err)
		return
	}

	response.SuccessWithMessage(c, "退出成功", nil)
}
