package service

import (
	dto "blog/internal/model/dto/response"
	"blog/internal/repository"
	"blog/pkg/errors"
)

// authorService 作者服务实现
type authorService struct {
	repo repository.AuthorRepository
}

// NewAuthorService 创建作者服务
func NewAuthorService(repo repository.AuthorRepository) AuthorService {
	return &authorService{repo: repo}
}

// GetAuthorDetail 获取作者详情
func (s *authorService) GetAuthorDetail() (*dto.AuthorDTO, error) {
	author, err := s.repo.GetAuthor()
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "查询作者信息失败")
	}

	articleCount, err := s.repo.CountArticles(author.ID)
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "统计文章数量失败")
	}

	categoryCount, err := s.repo.CountCategories(author.ID)
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "统计分类数量失败")
	}

	tagCount, err := s.repo.CountTags(author.ID)
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "统计标签数量失败")
	}

	return &dto.AuthorDTO{
		ID:            author.ID,
		Nickname:      author.Nickname,
		Bio:           author.Bio,
		Avatar:        author.Avatar,
		Github:        author.Github,
		ArticleCount:  articleCount,
		CategoryCount: categoryCount,
		TagCount:      tagCount,
	}, nil
}
