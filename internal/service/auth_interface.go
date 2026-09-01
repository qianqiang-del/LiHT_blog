package service

import (
	"blog/internal/model/dto/request"
	dto "blog/internal/model/dto/response"
)

// AuthService 认证服务接口
type AuthService interface {
	// SendCode 发送验证码
	SendCode(req request.SendCodeRequest) error
	// Register 用户注册（注册成功自动登录）
	Register(req request.RegisterRequest) (*dto.LoginResponse, error)
	// Login 用户登录
	Login(req request.LoginRequest) (*dto.LoginResponse, error)
	// Logout 退出登录（将 token 加入黑名单）
	Logout(token string) error
	// GenerateCaptcha 生成图形验证码
	GenerateCaptcha() (*dto.CaptchaResponse, error)
	// GetUserInfo 获取用户信息
	GetUserInfo(userID uint) (*dto.UserInfo, error)
}
