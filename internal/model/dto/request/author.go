package request

// AdminLoginRequest 后台登录请求
type AdminLoginRequest struct {
	Account     string `json:"account" binding:"required,max=50"`
	Password    string `json:"password" binding:"required,max=50"`
	CaptchaID   string `json:"captcha_id" binding:"required"`
	CaptchaCode string `json:"captcha_code" binding:"required"`
}

// AdminUpdateAccountRequest 修改登录账号请求
type AdminUpdateAccountRequest struct {
	Account string `json:"account" binding:"required,max=50"`
}

// AdminUpdatePasswordRequest 修改登录密码请求
type AdminUpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,max=50"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}

// AdminUpdateAuthorRequest 后台更新作者信息请求
type AdminUpdateAuthorRequest struct {
	Nickname   *string `json:"nickname"`
	Avatar     *string `json:"avatar"`
	Bio        *string `json:"bio"`
	Github     *string `json:"github"`
	About      *string `json:"about"`
	Background *string `json:"background"`
	AboutBlog  *string `json:"about_blog"`
}
