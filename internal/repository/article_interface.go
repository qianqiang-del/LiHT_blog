package repository

import (
	"blog/internal/model/entity"

	"gorm.io/gorm"
)

// ArticleRepository 文章仓储接口
type ArticleRepository interface {
	ListPublished(offset, limit int) ([]entity.Article, int64, error)
	ListHot(offset, limit int) ([]entity.Article, int64, error)
	// SearchByKeyword 全文搜索文章（标题、描述、正文）
	SearchByKeyword(keyword string, offset, limit int) ([]entity.Article, int64, error)
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
	// GetLikeCount 获取最新点赞数
	GetLikeCount(db *gorm.DB, articleID uint) (int, error)
	// IncrementCommentCount 文章评论数 +1
	IncrementCommentCount(db *gorm.DB, articleID uint) error
	// DecrementCommentCount 文章评论数 -1
	DecrementCommentCount(db *gorm.DB, articleID uint) error
	// AdminListArticles 后台查询文章列表（支持标题模糊查询和分类筛选）
	AdminListArticles(title, category string, offset, limit int) ([]entity.Article, int64, error)
	// UpdateStatus 更新文章状态
	UpdateStatus(id uint, status int8) error
	// HardDelete 硬删除文章
	HardDelete(id uint) error
	// Create 创建文章（含标签关联）
	Create(article *entity.Article, tagIDs []uint) error
}
