package service

import dto "blog/internal/model/dto/response"

// CategoryService 分类服务接口
type CategoryService interface {
	GetCategoryDetail() ([]dto.CategoryDTO, error)
}
