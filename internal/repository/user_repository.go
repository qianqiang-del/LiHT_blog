package repository

import (
	"blog/internal/model/entity"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Count 统计用户总数
func (r *userRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&entity.User{}).Count(&count).Error
	return count, err
}

// AdminList 后台获取用户列表
func (r *userRepository) AdminList(offset, limit int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	query := r.db.Model(&entity.User{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("id DESC").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetByID 根据 ID 获取用户
func (r *userRepository) GetByID(id uint) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateStatus 更新用户状态
func (r *userRepository) UpdateStatus(id uint, status int8) error {
	return r.db.Model(&entity.User{}).Where("id = ?", id).Update("status", status).Error
}

// Delete 硬删除用户（评论级联删除，点赞删除并更新计数）
func (r *userRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. 文章点赞：先减少文章的 like_count，再删除点赞记录
		if err := tx.Model(&entity.Article{}).
			Where("id IN (SELECT article_id FROM article_likes WHERE user_id = ?)", id).
			UpdateColumn("like_count", gorm.Expr("GREATEST(like_count - 1, 0)")).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", id).Delete(&entity.ArticleLike{}).Error; err != nil {
			return err
		}

		// 2. 评论点赞：先减少评论的 like_count，再删除点赞记录
		if err := tx.Model(&entity.Comment{}).
			Where("id IN (SELECT comment_id FROM comment_likes WHERE user_id = ?)", id).
			UpdateColumn("like_count", gorm.Expr("GREATEST(like_count - 1, 0)")).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", id).Delete(&entity.CommentLike{}).Error; err != nil {
			return err
		}

		// 3. 评论级联删除
		// 3a. 查询该用户的所有评论
		var userComments []entity.Comment
		if err := tx.Unscoped().Where("user_id = ?", id).Find(&userComments).Error; err != nil {
			return err
		}

		// 按文章统计需要减少的评论数
		articleDecr := make(map[uint]int)
		var userCommentIDs []uint
		var userL1IDs []uint

		for _, c := range userComments {
			userCommentIDs = append(userCommentIDs, c.ID)
			articleDecr[c.ArticleID]++
			if c.ParentID == nil {
				userL1IDs = append(userL1IDs, c.ID)
			}
		}

		// 3b. 删除用户一级评论下的所有子评论（包括其他用户的回复）
		if len(userL1IDs) > 0 {
			var childComments []entity.Comment
			if err := tx.Unscoped().Where("parent_id IN ?", userL1IDs).Find(&childComments).Error; err != nil {
				return err
			}
			for _, c := range childComments {
				articleDecr[c.ArticleID]++
			}
			if err := tx.Unscoped().Where("parent_id IN ?", userL1IDs).Delete(&entity.Comment{}).Error; err != nil {
				return err
			}
		}

		// 3c. 删除用户自己的评论
		if len(userCommentIDs) > 0 {
			if err := tx.Unscoped().Where("id IN ?", userCommentIDs).Delete(&entity.Comment{}).Error; err != nil {
				return err
			}
		}

		// 3d. 更新各文章的 comment_count
		for articleID, decr := range articleDecr {
			if err := tx.Model(&entity.Article{}).Where("id = ?", articleID).
				UpdateColumn("comment_count", gorm.Expr("GREATEST(comment_count - ?, 0)", decr)).Error; err != nil {
				return err
			}
		}

		// 4. 删除用户
		return tx.Unscoped().Delete(&entity.User{}, id).Error
	})
}

// CountComments 统计用户的评论数
func (r *userRepository) CountComments(userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&entity.Comment{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// BatchCountComments 批量统计用户评论数
func (r *userRepository) BatchCountComments(userIDs []uint) (map[uint]int64, error) {
	if len(userIDs) == 0 {
		return make(map[uint]int64), nil
	}

	type Result struct {
		UserID uint
		Count  int64
	}
	var results []Result

	if err := r.db.Model(&entity.Comment{}).
		Where("user_id IN ?", userIDs).
		Group("user_id").
		Select("user_id, COUNT(*) as count").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	counts := make(map[uint]int64, len(results))
	for _, result := range results {
		counts[result.UserID] = result.Count
	}
	return counts, nil
}

// SumLikeCount 统计用户点赞文章的数量
func (r *userRepository) SumLikeCount(userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&entity.ArticleLike{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// BatchSumLikeCount 批量统计用户点赞文章的数量
func (r *userRepository) BatchSumLikeCount(userIDs []uint) (map[uint]int64, error) {
	if len(userIDs) == 0 {
		return make(map[uint]int64), nil
	}

	type Result struct {
		UserID uint
		Count  int64
	}
	var results []Result

	if err := r.db.Model(&entity.ArticleLike{}).
		Where("user_id IN ?", userIDs).
		Group("user_id").
		Select("user_id, COUNT(*) as count").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	counts := make(map[uint]int64, len(results))
	for _, result := range results {
		counts[result.UserID] = result.Count
	}
	return counts, nil
}
