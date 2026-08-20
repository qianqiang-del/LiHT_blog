package service

import (
	"blog/internal/model/dto/request"
	"blog/pkg/response"
)

// ArticleService 文章服务接口
type ArticleService interface {
	ListArticles(req request.ArticleListRequest) (*response.PageResponse, error)
	ListHotArticles(req request.ArticleListRequest) (*response.PageResponse, error)
}
