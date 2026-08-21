package service

import (
	"blog/internal/model/dto/request"
	dto "blog/internal/model/dto/response"
	"blog/pkg/response"
)

// CategoryService 分类服务接口
type CategoryService interface {
	GetCategoryDetail() ([]dto.CategoryDTO, error)
	ListCategoryArticles(categoryID uint, req request.ArticleListRequest) (*response.PageResponse, error)
}
