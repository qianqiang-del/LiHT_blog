package repository

import (
	"blog/internal/model/entity"

	"gorm.io/gorm"
)

// commentRepository 评论仓储实现
type commentRepository struct {
	db *gorm.DB
}

// NewCommentRepository 创建评论仓储
func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

// Create 创建评论
func (r *commentRepository) Create(comment *entity.Comment) error {
	return r.db.Create(comment).Error
}

// ListByArticleID 获取文章的一级评论列表（分页，按时间倒序）
func (r *commentRepository) ListByArticleID(articleID uint, offset, limit int) ([]entity.Comment, int64, error) {
	var comments []entity.Comment
	var total int64

	query := r.db.Model(&entity.Comment{}).
		Where("article_id = ? AND parent_id IS NULL AND status = ?", articleID, 1)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// ListByParentID 获取评论的回复列表（分页，按时间正序）
func (r *commentRepository) ListByParentID(parentID uint, offset, limit int) ([]entity.Comment, int64, error) {
	var comments []entity.Comment
	var total int64

	query := r.db.Model(&entity.Comment{}).
		Where("parent_id = ? AND status = ?", parentID, 1)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("created_at ASC").
		Offset(offset).
		Limit(limit).
		Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// CountReplies 统计多条评论的回复数
func (r *commentRepository) CountReplies(commentIDs []uint) (map[uint]int, error) {
	if len(commentIDs) == 0 {
		return nil, nil
	}

	type result struct {
		ParentID uint
		Count    int
	}
	var results []result

	if err := r.db.Model(&entity.Comment{}).
		Select("parent_id, COUNT(*) as count").
		Where("parent_id IN ? AND status = ?", commentIDs, 1).
		Group("parent_id").
		Find(&results).Error; err != nil {
		return nil, err
	}

	countMap := make(map[uint]int, len(results))
	for _, r := range results {
		countMap[r.ParentID] = r.Count
	}
	return countMap, nil
}

// HasLiked 查询用户是否已点赞某评论
func (r *commentRepository) HasLiked(db *gorm.DB, commentID, userID uint) (bool, error) {
	var count int64
	err := db.Model(&entity.CommentLike{}).
		Where("comment_id = ? AND user_id = ?", commentID, userID).
		Count(&count).Error
	return count > 0, err
}

// CreateLike 创建评论点赞记录
func (r *commentRepository) CreateLike(db *gorm.DB, commentID, userID uint) error {
	return db.Create(&entity.CommentLike{CommentID: commentID, UserID: userID}).Error
}

// DeleteLike 删除评论点赞记录
func (r *commentRepository) DeleteLike(db *gorm.DB, commentID, userID uint) error {
	return db.Where("comment_id = ? AND user_id = ?", commentID, userID).
		Delete(&entity.CommentLike{}).Error
}

// IncrementLikeCount 评论点赞数 +1
func (r *commentRepository) IncrementLikeCount(db *gorm.DB, commentID uint) error {
	return db.Model(&entity.Comment{}).
		Where("id = ?", commentID).
		UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
}

// DecrementLikeCount 评论点赞数 -1
func (r *commentRepository) DecrementLikeCount(db *gorm.DB, commentID uint) error {
	return db.Model(&entity.Comment{}).
		Where("id = ?", commentID).
		UpdateColumn("like_count", gorm.Expr("like_count - 1")).Error
}

// GetLikeCount 获取评论最新点赞数
func (r *commentRepository) GetLikeCount(db *gorm.DB, commentID uint) (int, error) {
	var comment entity.Comment
	if err := db.Select("like_count").Where("id = ?", commentID).First(&comment).Error; err != nil {
		return 0, err
	}
	return comment.LikeCount, nil
}
