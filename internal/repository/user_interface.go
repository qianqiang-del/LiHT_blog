package repository

import "blog/internal/model/entity"

// UserRepository 用户仓储接口
type UserRepository interface {
	// AdminList 后台获取用户列表（分页）
	AdminList(offset, limit int) ([]entity.User, int64, error)
	// GetByID 根据 ID 获取用户
	GetByID(id uint) (*entity.User, error)
	// CountComments 统计用户的评论数
	CountComments(userID uint) (int64, error)
	// BatchCountComments 批量统计用户评论数
	BatchCountComments(userIDs []uint) (map[uint]int64, error)
	// SumLikeCount 统计用户获得的点赞总数（评论的点赞数之和）
	SumLikeCount(userID uint) (int64, error)
	// BatchSumLikeCount 批量统计用户获得的点赞总数
	BatchSumLikeCount(userIDs []uint) (map[uint]int64, error)
}
