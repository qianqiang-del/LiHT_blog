package repository

import (
	"blog/internal/model/entity"
	"gorm.io/gorm"
)

// categoryRepository 分类仓储实现
type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类仓储
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// ListCategories 获取分类列表
func (r *categoryRepository) ListCategories() ([]entity.Category, error) {
	var categories []entity.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

// CountArticles 统计分类下的已发布的文章数量
func (r *categoryRepository) CountArticles(categoryID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Article{}).
		Where("category_id = ? AND status = ?", categoryID, 1).
		Count(&count).Error
	return count, err
}

// ListByCategoryID 查询分类下的已发布文章（分页）
func (r *categoryRepository) ListByCategoryID(categoryID uint, offset, limit int) ([]entity.Article, int64, error) {
	var articles []entity.Article
	var total int64

	query := r.db.Model(&entity.Article{}).Where("category_id = ? AND status = ?", categoryID, 1)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("published_at DESC").Offset(offset).Limit(limit).Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}
