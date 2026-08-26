package repository

import "blog/internal/model/entity"

// AuthRepository 认证仓储接口
type AuthRepository interface {
	// FindByUsername 根据用户名查找用户
	FindByUsername(username string) (*entity.User, error)
	// FindByEmail 根据邮箱查找用户
	FindByEmail(email string) (*entity.User, error)
	// FindByAccount 根据用户名或邮箱查找用户
	FindByAccount(account string) (*entity.User, error)
	// CreateUser 创建用户
	CreateUser(user *entity.User) error
}
