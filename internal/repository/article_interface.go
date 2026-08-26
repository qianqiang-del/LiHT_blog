package repository

import (
	"blog/internal/model/entity"

	"gorm.io/gorm"
)

// ArticleRepository 文章仓储接口
type ArticleRepository interface {
	ListPublished(offset, limit int) ([]entity.Article, int64, error)
	ListHot(offset, limit int) ([]entity.Article, int64, error)
	// GetArticleTagNames 查询多篇文章的标签名称（article_id -> []tag_name）
	GetArticleTagNames(articleIDs []uint) (map[uint][]string, error)
	// GetArticleByID 根据 ID 获取文章详情（预加载分类）
	GetArticleByID(id uint) (*entity.Article, error)
	// GetArticleTags 获取文章的标签列表（id + name）
	GetArticleTags(articleID uint) ([]entity.Tag, error)
	// HasLiked 查询用户是否已点赞某文章
	HasLiked(db *gorm.DB, articleID, userID uint) (bool, error)
	// CreateLike 创建点赞记录
	CreateLike(db *gorm.DB, articleID, userID uint) error
	// DeleteLike 删除点赞记录
	DeleteLike(db *gorm.DB, articleID, userID uint) error
	// IncrementLikeCount 点赞数 +1
	IncrementLikeCount(db *gorm.DB, articleID uint) error
	// DecrementLikeCount 点赞数 -1
	DecrementLikeCount(db *gorm.DB, articleID uint) error
	// GetLikeCount 获取最新点赞数
	GetLikeCount(db *gorm.DB, articleID uint) (int, error)
}
