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

// Delete 硬删除用户
func (r *userRepository) Delete(id uint) error {
	return r.db.Unscoped().Delete(&entity.User{}, id).Error
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
