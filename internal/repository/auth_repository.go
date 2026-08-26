package repository

import (
	"blog/internal/model/entity"

	"gorm.io/gorm"
)

// authRepository 认证仓储实现
type authRepository struct {
	db *gorm.DB
}

// NewAuthRepository 创建认证仓储
func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

// FindByUsername 根据用户名查找用户
func (r *authRepository) FindByUsername(username string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail 根据邮箱查找用户
func (r *authRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByAccount 根据用户名或邮箱查找用户
func (r *authRepository) FindByAccount(account string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("username = ? OR email = ?", account, account).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser 创建用户
func (r *authRepository) CreateUser(user *entity.User) error {
	return r.db.Create(user).Error
}
