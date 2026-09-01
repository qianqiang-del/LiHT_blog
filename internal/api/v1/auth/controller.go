package auth

import (
	"strings"

	"blog/internal/middleware"
	"blog/internal/model/dto/request"
	"blog/internal/service"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
		response.BadRequest(c, parseValidationErr(err))
		return
	}

	result, err := ctrl.svc.Register(req)
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, result)
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

// parseValidationErr 将校验错误转为友好提示
func parseValidationErr(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			field := e.Field()
			switch field {
			case "Username":
				switch e.Tag() {
				case "min":
					return "用户名至少3个字符"
				case "max":
					return "用户名不能超过50个字符"
				}
			case "Password":
				switch e.Tag() {
				case "min":
					return "密码至少6个字符"
				}
			case "Email":
				return "邮箱格式不正确"
			case "Code":
				return "验证码格式不正确"
			}
		}
	}
	return "请求参数错误"
}
