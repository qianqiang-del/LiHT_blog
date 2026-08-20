package repository

import (
	"blog/internal/model/entity"

	"gorm.io/gorm"
)

// articleRepository 文章仓储实现
type articleRepository struct {
	db *gorm.DB
}

// NewArticleRepository 创建文章仓储
func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

// ListPublished 查询已发布的文章列表（按发布时间倒序）
func (r *articleRepository) ListPublished(offset, limit int) ([]entity.Article, int64, error) {
	var articles []entity.Article
	var total int64

	// 统计总数
	if err := r.db.Model(&entity.Article{}).
		Where("status = ?", 1).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询列表（按发布时间倒序）
	if err := r.db.
		Where("status = ?", 1).
		Order("published_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// ListHot 查询热门文章列表（按热门权重倒序）
func (r *articleRepository) ListHot(offset, limit int) ([]entity.Article, int64, error) {
	var articles []entity.Article
	var total int64

	// 统计总数
	if err := r.db.Model(&entity.Article{}).
		Where("status = ?", 1).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询列表（按热门权重倒序）
	if err := r.db.
		Where("status = ?", 1).
		Order("hot DESC").
		Offset(offset).
		Limit(limit).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// GetArticleTagNames 查询多篇文章的标签名称（article_id -> []tag_name）
func (r *articleRepository) GetArticleTagNames(articleIDs []uint) (map[uint][]string, error) {
	if len(articleIDs) == 0 {
		return nil, nil
	}

	var results []struct {
		ArticleID uint
		TagName   string
	}
	if err := r.db.Table("article_tags").
		Select("article_tags.article_id AS article_id, tags.name AS tag_name").
		Joins("JOIN tags ON tags.id = article_tags.tag_id").
		Where("article_tags.article_id IN ?", articleIDs).
		Scan(&results).Error; err != nil {
		return nil, err
	}

	tagMap := make(map[uint][]string, len(articleIDs))
	for _, r1 := range results {
		tagMap[r1.ArticleID] = append(tagMap[r1.ArticleID], r1.TagName)
	}
	return tagMap, nil
}
