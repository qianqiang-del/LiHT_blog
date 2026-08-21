package service

import (
	dto "blog/internal/model/dto/response"
	"blog/internal/repository"
	"blog/pkg/errors"
)

// tagService 标签服务实现
type tagService struct {
	repo repository.TagRepository
}

// NewTagService 创建标签服务
func NewTagService(repo repository.TagRepository) TagService {
	return &tagService{repo: repo}
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
