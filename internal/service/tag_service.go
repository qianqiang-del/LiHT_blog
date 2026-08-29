package service

import (
	"blog/internal/model/dto/request"
	dto "blog/internal/model/dto/response"
	"blog/internal/model/entity"
	"blog/internal/repository"
	"blog/pkg/errors"
	"blog/pkg/response"
)

// tagService 标签服务实现
type tagService struct {
	repo        repository.TagRepository
	articleRepo repository.ArticleRepository
}

// NewTagService 创建标签服务
func NewTagService(repo repository.TagRepository, articleRepo repository.ArticleRepository) TagService {
	return &tagService{repo: repo, articleRepo: articleRepo}
}

// ListTags 获取标签列表
func (s *tagService) ListTags() ([]dto.TagDTO, error) {
	tags, err := s.repo.ListTags()
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "查询标签列表失败")
	}

	list := make([]dto.TagDTO, 0, len(tags))
	for _, t := range tags {
		count, err := s.repo.CountArticles(t.ID)
		if err != nil {
			return nil, errors.New(errors.CodeInternalError, "统计文章数量失败")
		}
		list = append(list, dto.TagDTO{
			ID:           t.ID,
			Name:         t.Name,
			ArticleCount: count,
		})
	}

	return list, nil
}

// ListTagArticles 获取标签下的文章列表
func (s *tagService) ListTagArticles(tagIDs []uint, req request.ArticleListRequest) (*response.PageResponse, error) {
	offset := (req.Page - 1) * req.Size

	articles, total, err := s.repo.ListByTagIDs(tagIDs, offset, req.Size)
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "查询标签文章失败")
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

// AdminCreateTag 后台添加标签
func (s *tagService) AdminCreateTag(name string) error {
	tag := &entity.Tag{
		Name: name,
	}
	if err := s.repo.Create(tag); err != nil {
		return errors.New(errors.CodeInternalError, "添加标签失败")
	}
	return nil
}
