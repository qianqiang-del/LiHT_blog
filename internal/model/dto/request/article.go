package request

import "blog/pkg/response"

// ArticleListRequest 文章列表请求
type ArticleListRequest struct {
	response.PageRequest
}
