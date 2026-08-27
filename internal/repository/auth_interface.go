package repository

import (
	"blog/internal/model/entity"
	"time"
)

// AuthRepository 认证仓储接口
type AuthRepository interface {
	// FindByID 根据 ID 查询用户
	FindByID(id uint) (*entity.User, error)
	// FindByUsername 用户操作
	FindByUsername(username string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	FindByAccount(account string) (*entity.User, error)
	CreateUser(user *entity.User) error

	// SaveEmailCode 邮箱验证码
	SaveEmailCode(email, code string, ttl time.Duration) error
	GetEmailCode(email string) (string, error)
	DeleteEmailCode(email string) error

	// Set 图形验证码（实现 base64Captcha.Store 接口）
	Set(id string, value string) error
	Get(id string, clear bool) string
	Verify(id, answer string, clear bool) bool

	// AddTokenBlacklist Token 黑名单
	AddTokenBlacklist(token string, ttl time.Duration) error
	IsTokenBlacklisted(token string) (bool, error)
}
