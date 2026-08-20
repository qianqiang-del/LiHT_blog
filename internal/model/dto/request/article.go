package request

import "blog/pkg/response"

// ArticleListRequest 文章列表请求（仅 page + size）
type ArticleListRequest struct {
	response.PageRequest
}
