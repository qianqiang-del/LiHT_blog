package service

import (
	"blog/internal/model/dto/request"
	dto "blog/internal/model/dto/response"
	"blog/pkg/response"
)

// ArticleService 文章服务接口
type ArticleService interface {
	ListArticles(req request.ArticleListRequest) (*response.PageResponse, error)
	ListHotArticles(req request.ArticleListRequest) (*response.PageResponse, error)
	GetArticleDetail(id uint, userID *uint) (*dto.ArticleDetail, error)
	LikeArticle(articleID, userID uint) (*dto.LikeResponse, error)
}
