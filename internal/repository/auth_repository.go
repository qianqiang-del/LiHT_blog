package repository

import (
	"blog/internal/model/entity"
	"context"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	codePrefix      = "register:code:"   // 邮箱验证码 key 前缀
	captchaPrefix   = "captcha:"         // 图形验证码 key 前缀
	blacklistPrefix = "token:blacklist:" // token 黑名单 key 前缀
)

// authRepository 认证仓储实现
type authRepository struct {
	db    *gorm.DB
	redis RedisRepository
}

// NewAuthRepository 创建认证仓储
func NewAuthRepository(db *gorm.DB, redis RedisRepository) AuthRepository {
	return &authRepository{db: db, redis: redis}
}

// FindByID 根据 ID 查找用户
func (r *authRepository) FindByID(id uint) (*entity.User, error) {
	var user entity.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
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

// SaveEmailCode 保存邮箱验证码
func (r *authRepository) SaveEmailCode(email, code string, ttl time.Duration) error {
	return r.redis.Set(context.Background(), codePrefix+email, code, ttl)
}

// GetEmailCode 获取邮箱验证码
func (r *authRepository) GetEmailCode(email string) (string, error) {
	var code string
	err := r.redis.Get(context.Background(), codePrefix+email, &code)
	return code, err
}

// DeleteEmailCode 删除邮箱验证码
func (r *authRepository) DeleteEmailCode(email string) error {
	return r.redis.Del(context.Background(), codePrefix+email)
}

// Set 保存图形验证码（实现 base64Captcha.Store 接口）
func (r *authRepository) Set(id string, value string) error {
	return r.redis.Set(context.Background(), captchaPrefix+id, value, 5*time.Minute)
}

// Get 获取图形验证码（实现 base64Captcha.Store 接口）
func (r *authRepository) Get(id string, clear bool) string {
	var code string
	if err := r.redis.Get(context.Background(), captchaPrefix+id, &code); err != nil {
		return ""
	}
	if clear {
		_ = r.redis.Del(context.Background(), captchaPrefix+id)
	}
	return code
}

// Verify 校验图形验证码（实现 base64Captcha.Store 接口）
func (r *authRepository) Verify(id, answer string, clear bool) bool {
	code := r.Get(id, clear)
	if code == "" {
		return false
	}
	return strings.EqualFold(code, answer)
}

// AddTokenBlacklist 将 token 加入黑名单
func (r *authRepository) AddTokenBlacklist(token string, ttl time.Duration) error {
	return r.redis.Set(context.Background(), blacklistPrefix+token, "1", ttl)
}

// IsTokenBlacklisted 检查 token 是否在黑名单中
func (r *authRepository) IsTokenBlacklisted(token string) (bool, error) {
	return r.redis.Exists(context.Background(), blacklistPrefix+token)
}
