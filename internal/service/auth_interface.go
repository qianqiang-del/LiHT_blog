package service

import (
	"blog/internal/model/dto/request"
	dto "blog/internal/model/dto/response"
)

// AuthService 认证服务接口
type AuthService interface {
	// SendCode 发送验证码
	SendCode(req request.SendCodeRequest) error
	// Register 用户注册
	Register(req request.RegisterRequest) error
	// Login 用户登录
	Login(req request.LoginRequest) (*dto.LoginResponse, error)
	// Logout 退出登录（将 token 加入黑名单）
	Logout(token string) error
	// GenerateCaptcha 生成图形验证码
	GenerateCaptcha() (*dto.CaptchaResponse, error)
}
