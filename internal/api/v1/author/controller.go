package author

import (
	"strings"

	"blog/internal/middleware"
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

// adminLogin POST /api/v1/admin/login
func (ctrl *Controller) adminLogin(c *gin.Context) {
	var req request.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	result, err := ctrl.svc.AdminLogin(req)
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, result)
}

// adminLogout POST /api/v1/admin/logout
func (ctrl *Controller) adminLogout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		response.BadRequest(c, "令牌格式错误")
		return
	}

	if err := ctrl.svc.AdminLogout(parts[1]); err != nil {
		response.BizError(c, err)
		return
	}

	response.SuccessWithMessage(c, "退出成功", nil)
}

// adminGetCaptcha GET /api/v1/admin/captcha
func (ctrl *Controller) adminGetCaptcha(c *gin.Context) {
	result, err := ctrl.svc.AdminGetCaptcha()
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

// adminUpdateAccount PUT /api/v1/admin/author/account
func (ctrl *Controller) adminUpdateAccount(c *gin.Context) {
	authorID := middleware.GetUserID(c)
	if authorID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}

	var req request.AdminUpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	if err := ctrl.svc.AdminUpdateAccount(authorID, req); err != nil {
		response.BizError(c, err)
		return
	}

	response.SuccessWithMessage(c, "修改账号成功", nil)
}

// adminUpdatePassword PUT /api/v1/admin/author/password
func (ctrl *Controller) adminUpdatePassword(c *gin.Context) {
	authorID := middleware.GetUserID(c)
	if authorID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}

	var req request.AdminUpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	if err := ctrl.svc.AdminUpdatePassword(authorID, req); err != nil {
		response.BizError(c, err)
		return
	}

	response.SuccessWithMessage(c, "修改密码成功", nil)
}
