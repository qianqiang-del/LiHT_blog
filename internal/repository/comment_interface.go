package repository

import (
	"blog/internal/model/entity"

	"gorm.io/gorm"
)

// CommentRepository 评论仓储接口
type CommentRepository interface {
	// Create 创建评论
	Create(comment *entity.Comment) error
	// ListByArticleID 获取文章的一级评论列表（分页）
	ListByArticleID(articleID uint, offset, limit int) ([]entity.Comment, int64, error)
	// ListByParentID 获取评论的回复列表（分页）
	ListByParentID(parentID uint, offset, limit int) ([]entity.Comment, int64, error)
	// CountReplies 统计多条评论的回复数
	CountReplies(commentIDs []uint) (map[uint]int, error)
	// HasLiked 查询用户是否已点赞某评论
	HasLiked(db *gorm.DB, commentID, userID uint) (bool, error)
	// CreateLike 创建评论点赞记录
	CreateLike(db *gorm.DB, commentID, userID uint) error
	// DeleteLike 删除评论点赞记录
	DeleteLike(db *gorm.DB, commentID, userID uint) error
	// GetLikeCount 获取评论最新点赞数
	GetLikeCount(db *gorm.DB, commentID uint) (int, error)
}
