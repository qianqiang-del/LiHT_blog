package service

import (
	dto "blog/internal/model/dto/response"
	"blog/internal/repository"
	"blog/pkg/errors"
)

// categoryService 分类服务实现
type categoryService struct {
	repo repository.CategoryRepository
}

// NewCategoryService 创建分类服务
func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetCategoryDetail() ([]dto.CategoryDTO, error) {
	categories, err := s.repo.ListCategories()
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "查询分类标签表失败")
	}
	list := make([]dto.CategoryDTO, 0, len(categories))
	for _, c := range categories {
		count, err := s.repo.CountArticles(c.ID)
		if err != nil {
			return nil, errors.New(errors.CodeInternalError, "统计文章数量失败")
		}
		list = append(list, dto.CategoryDTO{
			ID:           c.ID,
			Name:         c.Name,
			ArticleCount: count,
		})
	}
	return list, nil
}
