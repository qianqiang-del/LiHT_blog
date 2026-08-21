package repository

import (
	"blog/internal/model/entity"

	"gorm.io/gorm"
)

// tagRepository 标签仓储实现
type tagRepository struct {
	db *gorm.DB
}

// NewTagRepository 创建标签仓储
func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

// ListTags 获取所有标签
func (r *tagRepository) ListTags() ([]entity.Tag, error) {
	var tags []entity.Tag
	err := r.db.Find(&tags).Error
	return tags, err
}

// CountArticles 统计标签下的已发布文章数量
func (r *tagRepository) CountArticles(tagID uint) (int64, error) {
	var count int64
	err := r.db.Table("article_tags").
		Joins("JOIN articles ON articles.id = article_tags.article_id").
		Where("article_tags.tag_id = ? AND articles.status = ?", tagID, 1).
		Count(&count).Error
	return count, err
}
