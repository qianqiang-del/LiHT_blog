package dto

// AdminLoginDTO 后台登录响应
type AdminLoginDTO struct {
	Token string `json:"token"`
}

// LoginResponse 前台登录响应
type LoginResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
	Email    string `json:"email"`
}

// CaptchaResponse 图形验证码响应
type CaptchaResponse struct {
	CaptchaID  string `json:"captcha_id"`
	CaptchaImg string `json:"captcha_img"`
}
