package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"blog/internal/model/dto/request"
	dto "blog/internal/model/dto/response"
	"blog/internal/model/entity"
	"blog/internal/repository"
	"blog/pkg/email"
	"blog/pkg/errors"
	"blog/pkg/jwt"

	"github.com/mojocn/base64Captcha"
	"golang.org/x/crypto/bcrypt"
)

const (
	codeTTL    = 1 * time.Minute // 邮箱验证码有效期
	codeLength = 6               // 邮箱验证码长度
)

// authService 认证服务实现
type authService struct {
	authRepo repository.AuthRepository
	email    *email.Sender
}

// NewAuthService 创建认证服务
func NewAuthService(authRepo repository.AuthRepository, email *email.Sender) AuthService {
	return &authService{
		authRepo: authRepo,
		email:    email,
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
	if err := s.authRepo.SaveEmailCode(req.Email, code, codeTTL); err != nil {
		return errors.New(errors.CodeInternalError, "存储验证码失败")
	}

	// 发送邮件
	if err := s.email.SendCode(req.Email, code); err != nil {
		fmt.Printf("邮件发送错误: %v\n", err)
		return errors.New(errors.CodeInternalError, "发送验证码失败: "+err.Error())
	}

	return nil
}

// GenerateCaptcha 生成图形验证码
func (s *authService) GenerateCaptcha() (*dto.CaptchaResponse, error) {
	driver := base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80)
	c := base64Captcha.NewCaptcha(driver, s.authRepo)
	id, b64s, _, err := c.Generate()
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "生成验证码失败")
	}
	return &dto.CaptchaResponse{
		CaptchaID:  id,
		CaptchaImg: b64s,
	}, nil
}

// verifyCaptcha 校验图形验证码
func (s *authService) verifyCaptcha(captchaID, captchaCode string) error {
	if !s.authRepo.Verify(captchaID, captchaCode, true) {
		return errors.New(errors.CodeInvalidParam, "验证码错误")
	}
	return nil
}

// Login 用户登录
func (s *authService) Login(req request.LoginRequest) (*dto.LoginResponse, error) {
	// 0. 校验图形验证码
	if err := s.verifyCaptcha(req.CaptchaID, req.CaptchaCode); err != nil {
		return nil, err
	}

	// 1. 根据用户名或邮箱查询用户
	user, err := s.authRepo.FindByAccount(req.Account)
	if err != nil {
		return nil, errors.New(errors.CodeInvalidCredentials, "用户名或密码错误")
	}

	// 2. 检查用户状态
	if user.Status != 1 {
		return nil, errors.New(errors.CodeUserDisabled, "账号已被禁用")
	}

	// 3. bcrypt 校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New(errors.CodeInvalidCredentials, "用户名或密码错误")
	}

	// 4. 生成 JWT token
	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "生成token失败")
	}

	// 5. 返回响应
	return &dto.LoginResponse{
		Token: token,
		User: dto.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}, nil
}

// Logout 退出登录，将 token 加入 Redis 黑名单
func (s *authService) Logout(token string) error {
	// 解析 token 获取过期时间
	claims, err := jwt.ParseToken(token)
	if err != nil {
		// token 已经无效，直接返回成功
		return nil
	}

	// 计算剩余过期时间
	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		return nil
	}

	// 存入黑名单
	return s.authRepo.AddTokenBlacklist(token, ttl)
}

// Register 用户注册
func (s *authService) Register(req request.RegisterRequest) error {
	// 验证码校验
	storedCode, err := s.authRepo.GetEmailCode(req.Email)
	if err != nil {
		return errors.New(errors.CodeInvalidParam, "验证码已过期或无效")
	}
	if storedCode != req.Code {
		return errors.New(errors.CodeInvalidParam, "验证码错误")
	}

	// 检查用户名是否已存在
	_, err = s.authRepo.FindByUsername(req.Username)
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
	_ = s.authRepo.DeleteEmailCode(req.Email)

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
	bytes := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
