package service

import (
	"blog/internal/model/dto/request"
	dto "blog/internal/model/dto/response"
	"blog/internal/model/entity"
	"blog/internal/repository"
	"blog/internal/stream"
	"blog/pkg/errors"
	"blog/pkg/response"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// service 文章服务实现
type service struct {
	repo      repository.ArticleRepository
	redisRepo repository.RedisRepository
	db        *gorm.DB
}

// NewArticleService 创建文章服务
func NewArticleService(repo repository.ArticleRepository, redisRepo repository.RedisRepository, db *gorm.DB) ArticleService {
	return &service{repo: repo, redisRepo: redisRepo, db: db}
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

// SearchArticles 全文搜索文章
func (s *service) SearchArticles(req request.ArticleSearchRequest) (*response.PageResponse, error) {
	offset := (req.Page - 1) * req.Size

	articles, total, err := s.repo.SearchByKeyword(req.Keyword, offset, req.Size)
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "搜索失败")
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
func (s *service) GetArticleDetail(id uint, userID *uint) (*dto.ArticleDetail, error) {
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

	// 查询点赞状态
	liked := false
	if userID != nil {
		liked, _ = s.repo.HasLiked(s.db, id, *userID)
	}

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
		Summary:      article.Summary,
		Content:      article.Content,
		Cover:        article.Cover,
		Category:     categoryDTO,
		Tags:         tagDTOs,
		ViewCount:    article.ViewCount + 1,
		LikeCount:    article.LikeCount,
		CommentCount: article.CommentCount,
		PublishedAt:  article.PublishedAt,
		Liked:        liked,
	}, nil
}

// LikeArticle 切换点赞状态
func (s *service) LikeArticle(articleID, userID uint) (*dto.LikeResponse, error) {
	var liked bool
	var likeCount int

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 查当前状态
		exists, err := s.repo.HasLiked(tx, articleID, userID)
		if err != nil {
			return err
		}

		if exists {
			// 已点赞 → 取消（DeleteLike 内部已处理 like_count -1）
			if err := s.repo.DeleteLike(tx, articleID, userID); err != nil {
				return err
			}
			liked = false
		} else {
			// 未点赞 → 点赞（CreateLike 内部已处理 like_count +1）
			if err := s.repo.CreateLike(tx, articleID, userID); err != nil {
				return err
			}
			liked = true
		}

		// 查最新点赞数
		count, err := s.repo.GetLikeCount(tx, articleID)
		if err != nil {
			return err
		}
		likeCount = count
		return nil
	})

	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "操作失败")
	}

	return &dto.LikeResponse{Liked: liked, LikeCount: likeCount}, nil
}

// nonNilTags 保证空标签返回 [] 而不是 null
func nonNilTags(tags []string) []string {
	if tags == nil {
		return []string{}
	}
	return tags
}

// AdminListArticles 后台获取文章列表
func (s *service) AdminListArticles(req request.AdminArticleListRequest) (*response.PageResponse, error) {
	offset := (req.Page - 1) * req.Size

	articles, total, err := s.repo.AdminListArticles(req.Title, req.Category, offset, req.Size)
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
	list := make([]dto.AdminArticleListItem, 0, len(articles))
	for _, a := range articles {
		// 获取分类名称
		categoryName := ""
		if a.Category != nil {
			categoryName = a.Category.Name
		}

		item := dto.AdminArticleListItem{
			ID:           a.ID,
			Title:        a.Title,
			Cover:        a.Cover,
			Category:     categoryName,
			Tags:         nonNilTags(tagMap[a.ID]),
			Status:       a.Status,
			ViewCount:    a.ViewCount,
			LikeCount:    a.LikeCount,
			CommentCount: a.CommentCount,
			Hot:          a.Hot,
			PublishedAt:  a.PublishedAt,
		}
		list = append(list, item)
	}

	return response.NewPageResponse(list, total, req.Page, req.Size), nil
}

// AdminUpdateArticleStatus 后台更新文章状态
func (s *service) AdminUpdateArticleStatus(id uint, status int8) error {
	// 检查文章是否存在
	_, err := s.repo.GetArticleByID(id)
	if err != nil {
		return errors.New(errors.CodeResourceNotFound, "文章不存在")
	}

	if err := s.repo.UpdateStatus(id, status); err != nil {
		return errors.New(errors.CodeInternalError, "更新文章状态失败")
	}

	return nil
}

// AdminDeleteArticle 后台硬删除文章
func (s *service) AdminDeleteArticle(id uint) error {
	// 检查文章是否存在
	_, err := s.repo.GetArticleByID(id)
	if err != nil {
		return errors.New(errors.CodeResourceNotFound, "文章不存在")
	}

	if err := s.repo.HardDelete(id); err != nil {
		return errors.New(errors.CodeInternalError, "删除文章失败")
	}

	return nil
}

// AdminCreateArticle 后台发布文章
func (s *service) AdminCreateArticle(req request.AdminCreateArticleRequest) error {
	now := time.Now()
	article := &entity.Article{
		Title:       req.Title,
		Summary:     req.Summary,
		Content:     req.Content,
		Cover:       req.Cover,
		AuthorID:    1, //
		CategoryID:  req.CategoryID,
		Status:      1, // 默认正常状态
		PublishedAt: &now,
	}

	if err := s.repo.Create(article, req.TagIDs); err != nil {
		return errors.New(errors.CodeInternalError, "发布文章失败")
	}

	return nil
}

// AdminUpdateArticle 后台更新文章
func (s *service) AdminUpdateArticle(id uint, req request.AdminUpdateArticleRequest) error {
	// 检查文章是否存在
	existing, err := s.repo.GetArticleByID(id)
	if err != nil {
		return errors.New(errors.CodeResourceNotFound, "文章不存在")
	}

	// 构建更新对象，只更新传入的字段
	article := &entity.Article{}
	article.ID = id

	if req.Title != nil {
		article.Title = *req.Title
	} else {
		article.Title = existing.Title
	}

	if req.Summary != nil {
		article.Summary = *req.Summary
	} else {
		article.Summary = existing.Summary
	}

	if req.Content != nil {
		article.Content = *req.Content
	} else {
		article.Content = existing.Content
	}

	if req.Cover != nil {
		article.Cover = *req.Cover
	} else {
		article.Cover = existing.Cover
	}

	if req.CategoryID != nil {
		article.CategoryID = *req.CategoryID
	} else {
		article.CategoryID = existing.CategoryID
	}

	// 标签：如果传了就用新的，否则用旧的
	tagIDs := req.TagIDs
	if tagIDs == nil {
		// 获取旧标签
		oldTags, err := s.repo.GetArticleTags(id)
		if err != nil {
			return errors.New(errors.CodeInternalError, "获取文章标签失败")
		}
		tagIDs = make([]uint, 0, len(oldTags))
		for _, t := range oldTags {
			tagIDs = append(tagIDs, t.ID)
		}
	}

	if err := s.repo.Update(article, tagIDs); err != nil {
		return errors.New(errors.CodeInternalError, "更新文章失败")
	}

	return nil
}

// AdminListArticleOptions 获取文章选项列表
func (s *service) AdminListArticleOptions() ([]dto.ArticleOption, error) {
	articles, err := s.repo.ListArticleOptions()
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "查询文章列表失败")
	}

	list := make([]dto.ArticleOption, 0, len(articles))
	for _, a := range articles {
		list = append(list, dto.ArticleOption{
			ID:    a.ID,
			Title: a.Title,
		})
	}

	return list, nil
}
