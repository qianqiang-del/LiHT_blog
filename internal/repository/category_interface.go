package repository

import "blog/internal/model/entity"

// CategoryRepository 分类仓储接口
type CategoryRepository interface {
	ListCategories() ([]entity.Category, error)
	CountArticles(categoryID uint) (int64, error)
	ListByCategoryID(categoryID uint, offset, limit int) ([]entity.Article, int64, error)
}
