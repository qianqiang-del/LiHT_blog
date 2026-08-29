package request

// AdminCreateCategoryRequest 后台添加分类请求
type AdminCreateCategoryRequest struct {
	Name string `json:"name" binding:"required,min=1,max=20"`
}
