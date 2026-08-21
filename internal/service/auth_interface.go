package service

import "blog/internal/model/dto/request"

// AuthService 认证服务接口
type AuthService interface {
	// SendCode 发送验证码
	SendCode(req request.SendCodeRequest) error
	// Register 用户注册
	Register(req request.RegisterRequest) error
}
