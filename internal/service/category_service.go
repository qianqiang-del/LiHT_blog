package service

import (
	"blog/internal/model/dto/request"
	dto "blog/internal/model/dto/response"
	"blog/internal/repository"
	"blog/pkg/errors"
	"blog/pkg/response"
)

// categoryService 分类服务实现
type categoryService struct {
	repo        repository.CategoryRepository
	articleRepo repository.ArticleRepository
}

// NewCategoryService 创建分类服务
func NewCategoryService(repo repository.CategoryRepository, articleRepo repository.ArticleRepository) CategoryService {
	return &categoryService{repo: repo, articleRepo: articleRepo}
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

// ListCategoryArticles 获取分类下的文章列表
func (s *categoryService) ListCategoryArticles(categoryID uint, req request.ArticleListRequest) (*response.PageResponse, error) {
	offset := (req.Page - 1) * req.Size

	articles, total, err := s.repo.ListByCategoryID(categoryID, offset, req.Size)
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "查询分类文章失败")
	}

	// 查询标签
	articleIDs := make([]uint, 0, len(articles))
	for _, a := range articles {
		articleIDs = append(articleIDs, a.ID)
	}
	tagMap, err := s.articleRepo.GetArticleTagNames(articleIDs)
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "查询文章标签失败")
	}

	// 组装 DTO
	list := make([]dto.ArticleListItem, 0, len(articles))
	for _, a := range articles {
		list = append(list, dto.ArticleListItem{
			ID:          a.ID,
			Title:       a.Title,
			Summary:     a.Summary,
			Cover:       a.Cover,
			Tags:        nonNilTags(tagMap[a.ID]),
			ViewCount:   a.ViewCount,
			PublishedAt: a.PublishedAt,
		})
	}

	return response.NewPageResponse(list, total, req.Page, req.Size), nil
}
