package service

import (
	"blog/internal/model/dto/request"
	dto "blog/internal/model/dto/response"
	"blog/pkg/response"
)

// TagService 标签服务接口
type TagService interface {
	ListTags() ([]dto.TagDTO, error)
	ListTagArticles(tagIDs []uint, req request.ArticleListRequest) (*response.PageResponse, error)
}
