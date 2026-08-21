package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"blog/internal/model/dto/request"
	"blog/internal/model/entity"
	"blog/internal/repository"
	"blog/pkg/email"
	"blog/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	codePrefix = "register:code:" // Redis key 前缀
	codeTTL    = 1 * time.Minute  // 验证码有效期
	codeLength = 6                // 验证码长度
)

// authService 认证服务实现
type authService struct {
	authRepo  repository.AuthRepository
	redisRepo repository.RedisRepository
	email     *email.Sender
}

// NewAuthService 创建认证服务
func NewAuthService(authRepo repository.AuthRepository, redisRepo repository.RedisRepository, email *email.Sender) AuthService {
	return &authService{
		authRepo:  authRepo,
		redisRepo: redisRepo,
		email:     email,
	}
}

// SendCode 发送验证码
func (s *authService) SendCode(req request.SendCodeRequest) error {
	// 检查邮箱是否已注册
	_, err := s.authRepo.FindByEmail(req.Email)
	if err == nil {
		return errors.New(errors.CodeConflict, "该邮箱已注册")
	}

	// 生成验证码
	code, err := generateCode()
	if err != nil {
		return errors.New(errors.CodeInternalError, "生成验证码失败")
	}

	// 存储到 Redis
	key := codePrefix + req.Email
	if err := s.redisRepo.Set(context.Background(), key, code, codeTTL); err != nil {
		return errors.New(errors.CodeInternalError, "存储验证码失败")
	}

	// 发送邮件
	if err := s.email.SendCode(req.Email, code); err != nil {
		fmt.Printf("邮件发送错误: %v\n", err)
		return errors.New(errors.CodeInternalError, "发送验证码失败: "+err.Error())
	}

	return nil
}

// Register 用户注册
func (s *authService) Register(req request.RegisterRequest) error {
	// 验证码校验
	key := codePrefix + req.Email
	var storedCode string
	if err := s.redisRepo.Get(context.Background(), key, &storedCode); err != nil {
		return errors.New(errors.CodeInvalidParam, "验证码已过期或无效")
	}
	if storedCode != req.Code {
		return errors.New(errors.CodeInvalidParam, "验证码错误")
	}

	// 检查用户名是否已存在
	_, err := s.authRepo.FindByUsername(req.Username)
	if err == nil {
		return errors.New(errors.CodeConflict, "用户名已存在")
	}

	// 检查邮箱是否已注册
	_, err = s.authRepo.FindByEmail(req.Email)
	if err == nil {
		return errors.New(errors.CodeConflict, "该邮箱已注册")
	}

	// 密码加密
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return errors.New(errors.CodeInternalError, "密码加密失败")
	}

	// 创建用户
	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Status:   1,
	}
	if err := s.authRepo.CreateUser(user); err != nil {
		return errors.New(errors.CodeInternalError, "创建用户失败")
	}

	// 删除已使用的验证码
	_ = s.redisRepo.Del(context.Background(), key)

	return nil
}

// generateCode 生成 6 位随机验证码
func generateCode() (string, error) {
	code := ""
	for i := 0; i < codeLength; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code += fmt.Sprintf("%d", n.Int64())
	}
	return code, nil
}

// hashPassword 密码加密
func hashPassword(password string) (string, error) {
	// 使用 bcrypt 加密
	bytes := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
