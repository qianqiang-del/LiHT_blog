package dto

// LoginResponse 登录响应
type LoginResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// CaptchaResponse 图形验证码响应
type CaptchaResponse struct {
	CaptchaID  string `json:"captcha_id"`
	CaptchaImg string `json:"captcha_image"` // base64 编码的 PNG 图片
}
