package repository

import "blog/internal/model/entity"

// TagRepository 标签仓储接口
type TagRepository interface {
	ListTags() ([]entity.Tag, error)
	CountArticles(tagID uint) (int64, error)
	ListByTagIDs(tagIDs []uint, offset, limit int) ([]entity.Article, int64, error)
	Create(tag *entity.Tag) error
}
