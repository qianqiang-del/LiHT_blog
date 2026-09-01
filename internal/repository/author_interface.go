package repository

import "blog/internal/model/entity"

// AuthorRepository 作者仓储接口
type AuthorRepository interface {
	// GetAuthor 获取唯一作者
	GetAuthor() (*entity.Author, error)
	// FindByID 根据 ID 查询作者
	FindByID(id uint) (*entity.Author, error)
	// FindByAccount 根据账号查询作者
	FindByAccount(account string) (*entity.Author, error)
	// Update 更新作者信息
	Update(author *entity.Author) error
	// CountArticles 统计作者的文章数量
	CountArticles(authorID uint) (int64, error)
	// CountCategories 统计作者的分类数量（去重）
	CountCategories(authorID uint) (int64, error)
	// CountTags 统计作者的标签数量（去重）
	CountTags(authorID uint) (int64, error)
}
