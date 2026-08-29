package request

import "blog/pkg/response"

// ArticleListRequest 文章列表请求
type ArticleListRequest struct {
	response.PageRequest
}

// ArticleSearchRequest 文章搜索请求
type ArticleSearchRequest struct {
	response.PageRequest
	Keyword string `form:"keyword" binding:"required,min=2,max=200"`
}

// AdminArticleListRequest 后台文章列表请求
type AdminArticleListRequest struct {
	response.PageRequest
	Title    string `form:"title"`    // 文章标题模糊查询
	Category string `form:"category"` // 分类筛选
}

// AdminUpdateArticleStatusRequest 后台更新文章状态请求
type AdminUpdateArticleStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=正常 下架"` // 正常=1 下架=0
}

// AdminCreateArticleRequest 后台发布文章请求
type AdminCreateArticleRequest struct {
	Title      string `json:"title" binding:"required,max=50"`
	Summary    string `json:"summary" binding:"max=200"`
	Content    string `json:"content" binding:"required,max=5000"`
	Cover      string `json:"cover"`
	CategoryID uint   `json:"category_id" binding:"required"`
	TagIDs     []uint `json:"tag_ids" binding:"max=3"`
}
