package request

// AdminUserListRequest 后台用户列表请求
type AdminUserListRequest struct {
	Page int `form:"page" binding:"required,min=1"`
	Size int `form:"size" binding:"required,min=1,max=50"`
}
