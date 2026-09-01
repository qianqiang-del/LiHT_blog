package repository

import (
	"blog/internal/model/entity"

	"gorm.io/gorm"
)

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

// Count 统计评论总数
func (r *commentRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&entity.Comment{}).Count(&count).Error
	return count, err
}

// Create 创建评论
func (r *commentRepository) Create(comment *entity.Comment) error {
	return r.db.Create(comment).Error
}

// ListByArticleID 获取文章的一级评论列表
func (r *commentRepository) ListByArticleID(articleID uint, offset, limit int) ([]entity.Comment, int64, error) {
	var comments []entity.Comment
	var total int64

	query := r.db.Model(&entity.Comment{}).
		Where("article_id = ? AND parent_id IS NULL", articleID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// ListByParentID 获取评论的回复列表
func (r *commentRepository) ListByParentID(parentID uint, offset, limit int) ([]entity.Comment, int64, error) {
	var comments []entity.Comment
	var total int64

	query := r.db.Model(&entity.Comment{}).
		Where("parent_id = ?", parentID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at ASC").Offset(offset).Limit(limit).Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// CountReplies 统计多条评论的回复数
func (r *commentRepository) CountReplies(commentIDs []uint) (map[uint]int, error) {
	if len(commentIDs) == 0 {
		return make(map[uint]int), nil
	}

	type Result struct {
		ParentID uint
		Count    int
	}
	var results []Result

	if err := r.db.Model(&entity.Comment{}).
		Where("parent_id IN ?", commentIDs).
		Group("parent_id").
		Select("parent_id, COUNT(*) as count").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	counts := make(map[uint]int, len(results))
	for _, result := range results {
		counts[result.ParentID] = result.Count
	}
	return counts, nil
}

// HasLiked 查询用户是否已点赞某评论
func (r *commentRepository) HasLiked(db *gorm.DB, commentID, userID uint) (bool, error) {
	var count int64
	if err := db.Table("comment_likes").Where("comment_id = ? AND user_id = ?", commentID, userID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// CreateLike 创建评论点赞记录
func (r *commentRepository) CreateLike(db *gorm.DB, commentID, userID uint) error {
	if err := db.Exec("INSERT INTO comment_likes (comment_id, user_id) VALUES (?, ?)", commentID, userID).Error; err != nil {
		return err
	}
	return db.Model(&entity.Comment{}).Where("id = ?", commentID).
		UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
}

// DeleteLike 删除评论点赞记录
func (r *commentRepository) DeleteLike(db *gorm.DB, commentID, userID uint) error {
	if err := db.Exec("DELETE FROM comment_likes WHERE comment_id = ? AND user_id = ?", commentID, userID).Error; err != nil {
		return err
	}
	return db.Model(&entity.Comment{}).Where("id = ?", commentID).
		UpdateColumn("like_count", gorm.Expr("GREATEST(like_count - 1, 0)")).Error
}

// GetLikeCount 获取评论最新点赞数
func (r *commentRepository) GetLikeCount(db *gorm.DB, commentID uint) (int, error) {
	var comment entity.Comment
	if err := db.Select("like_count").Where("id = ?", commentID).First(&comment).Error; err != nil {
		return 0, err
	}
	return comment.LikeCount, nil
}

// AdminList 后台查询评论列表
func (r *commentRepository) AdminList(offset, limit int, articleID *uint) ([]entity.Comment, int64, error) {
	var comments []entity.Comment
	var total int64
	query := r.db.Model(&entity.Comment{})
	if articleID != nil {
		query = query.Where("article_id = ?", *articleID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&comments).Error; err != nil {
		return nil, 0, err
	}
	return comments, total, nil
}

// Delete 后台硬删除评论（一级评论连带子评论一起删除，并更新文章评论数）
func (r *commentRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 先查询该评论，判断是否为一级评论
		var comment entity.Comment
		if err := tx.Unscoped().First(&comment, id).Error; err != nil {
			return err
		}

		deleteCount := 1 // 至少删除 1 条

		// 如果是一级评论（parent_id 为空），先删除所有子评论
		if comment.ParentID == nil {
			var childCount int64
			if err := tx.Unscoped().Model(&entity.Comment{}).Where("parent_id = ?", id).Count(&childCount).Error; err != nil {
				return err
			}
			if err := tx.Unscoped().Where("parent_id = ?", id).Delete(&entity.Comment{}).Error; err != nil {
				return err
			}
			deleteCount += int(childCount)
		}

		// 删除该评论本身
		if err := tx.Unscoped().Delete(&entity.Comment{}, id).Error; err != nil {
			return err
		}

		// 更新文章评论数
		return tx.Model(&entity.Article{}).Where("id = ?", comment.ArticleID).
			UpdateColumn("comment_count", gorm.Expr("GREATEST(comment_count - ?, 0)", deleteCount)).Error
	})
}
