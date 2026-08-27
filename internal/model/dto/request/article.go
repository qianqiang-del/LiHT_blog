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
