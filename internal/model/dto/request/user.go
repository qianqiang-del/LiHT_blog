package request

// AdminUserListRequest 后台用户列表请求
type AdminUserListRequest struct {
	Page int `form:"page" binding:"required,min=1"`
	Size int `form:"size" binding:"required,min=1,max=50"`
}

// AdminUpdateUserStatusRequest 后台更新用户状态请求
type AdminUpdateUserStatusRequest struct {
	Status *int8 `json:"status" binding:"required"` // 1=正常 0=禁用
}
