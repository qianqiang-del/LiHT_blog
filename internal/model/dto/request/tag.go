package request

// AdminCreateTagRequest 后台添加标签请求
type AdminCreateTagRequest struct {
	Name string `json:"name" binding:"required,min=1,max=20"`
}
