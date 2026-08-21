package service

import (
	"blog/internal/model/dto/request"
	dto "blog/internal/model/dto/response"
	"blog/internal/repository"
	"blog/internal/stream"
	"blog/pkg/errors"
	"blog/pkg/response"
	"context"
	"fmt"
)

// service 文章服务实现
type service struct {
	repo      repository.ArticleRepository
	redisRepo repository.RedisRepository
}

// NewArticleService 创建文章服务
func NewArticleService(repo repository.ArticleRepository, redisRepo repository.RedisRepository) ArticleService {
	return &service{repo: repo, redisRepo: redisRepo}
}

// ListArticles 获取文章列表
func (s *service) ListArticles(req request.ArticleListRequest) (*response.PageResponse, error) {
	offset := (req.Page - 1) * req.Size

	articles, total, err := s.repo.ListPublished(offset, req.Size)
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "查询文章列表失败")
	}

	// 一次查完所有文章的标签名称
	articleIDs := make([]uint, 0, len(articles))
	for _, a := range articles {
		articleIDs = append(articleIDs, a.ID)
	}
	tagMap, err := s.repo.GetArticleTagNames(articleIDs)
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "查询文章标签失败")
	}

	// 组装 DTO
	list := make([]dto.ArticleListItem, 0, len(articles))
	for _, a := range articles {
		item := dto.ArticleListItem{
			ID:          a.ID,
			Title:       a.Title,
			Summary:     a.Summary,
			Cover:       a.Cover,
			ViewCount:   a.ViewCount,
			PublishedAt: a.PublishedAt,
			Tags:        nonNilTags(tagMap[a.ID]),
		}
		list = append(list, item)
	}

	return response.NewPageResponse(list, total, req.Page, req.Size), nil
}

// ListHotArticles 获取热门文章列表
func (s *service) ListHotArticles(req request.ArticleListRequest) (*response.PageResponse, error) {
	offset := (req.Page - 1) * req.Size

	articles, total, err := s.repo.ListHot(offset, req.Size)
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "查询热门文章失败")
	}

	// 组装 DTO
	list := make([]dto.HotArticleItem, 0, len(articles))
	for _, a := range articles {
		item := dto.HotArticleItem{
			ID:          a.ID,
			Title:       a.Title,
			Summary:     a.Summary,
			Cover:       a.Cover,
			PublishedAt: a.PublishedAt,
		}
		list = append(list, item)
	}

	return response.NewPageResponse(list, total, req.Page, req.Size), nil
}

// GetArticleDetail 获取文章详情
func (s *service) GetArticleDetail(id uint) (*dto.ArticleDetail, error) {
	article, err := s.repo.GetArticleByID(id)
	if err != nil {
		return nil, errors.New(errors.CodeResourceNotFound, "文章不存在")
	}

	// 获取标签
	tags, err := s.repo.GetArticleTags(id)
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "查询文章标签失败")
	}

	// 浏览量 +1：发到 Redis Stream，由消费者异步写入 MySQL
	_ = s.redisRepo.XAdd(context.Background(), stream.StreamKey, map[string]interface{}{
		"id": fmt.Sprintf("%d", id),
	})

	// 组装标签 DTO
	tagDTOs := make([]dto.TagDTO, 0, len(tags))
	for _, t := range tags {
		tagDTOs = append(tagDTOs, dto.TagDTO{ID: t.ID, Name: t.Name})
	}

	// 组装分类 DTO
	var categoryDTO *dto.CategoryDTO
	if article.Category != nil {
		categoryDTO = &dto.CategoryDTO{ID: article.Category.ID, Name: article.Category.Name}
	}

	return &dto.ArticleDetail{
		ID:           article.ID,
		Title:        article.Title,
		Content:      article.Content,
		Cover:        article.Cover,
		Category:     categoryDTO,
		Tags:         tagDTOs,
		ViewCount:    article.ViewCount + 1, // 返回时包含本次浏览
		LikeCount:    article.LikeCount,
		CommentCount: article.CommentCount,
		PublishedAt:  article.PublishedAt,
		Liked:        false, // 默认未点赞
	}, nil
}

// nonNilTags 保证空标签返回 [] 而不是 null
func nonNilTags(tags []string) []string {
	if tags == nil {
		return []string{}
	}
	return tags
}
