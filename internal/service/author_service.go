package service

import (
	"time"

	"blog/internal/model/dto/request"
	dto "blog/internal/model/dto/response"
	"blog/internal/repository"
	"blog/pkg/errors"
	"blog/pkg/jwt"

	"github.com/mojocn/base64Captcha"
	"golang.org/x/crypto/bcrypt"
)

// authorService 作者服务实现
type authorService struct {
	repo     repository.AuthorRepository
	authRepo repository.AuthRepository
}

// NewAuthorService 创建作者服务
func NewAuthorService(repo repository.AuthorRepository, authRepo repository.AuthRepository) AuthorService {
	return &authorService{repo: repo, authRepo: authRepo}
}

// AdminLogin 后台登录
func (s *authorService) AdminLogin(req request.AdminLoginRequest) (*dto.AdminLoginDTO, error) {
	// 校验图形验证码
	if !s.authRepo.Verify(req.CaptchaID, req.CaptchaCode, true) {
		return nil, errors.New(errors.CodeInvalidParam, "验证码错误")
	}

	// 根据账号查询作者
	author, err := s.repo.FindByAccount(req.Account)
	if err != nil {
		return nil, errors.New(errors.CodeUnauthorized, "账号或密码错误")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(author.Password), []byte(req.Password)); err != nil {
		return nil, errors.New(errors.CodeUnauthorized, "账号或密码错误")
	}

	// 生成 JWT Token
	token, err := jwt.GenerateToken(author.ID, author.Account)
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "生成令牌失败")
	}

	return &dto.AdminLoginDTO{Token: token}, nil
}

// AdminLogout 后台退出登录，将 token 加入 Redis 黑名单
func (s *authorService) AdminLogout(token string) error {
	claims, err := jwt.ParseToken(token)
	if err != nil {
		return nil
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		return nil
	}

	return s.authRepo.AddTokenBlacklist(token, ttl)
}

// AdminGetCaptcha 获取图形验证码
func (s *authorService) AdminGetCaptcha() (*dto.CaptchaResponse, error) {
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

// GetAuthorDetail 获取作者详情
func (s *authorService) GetAuthorDetail() (*dto.AuthorDTO, error) {
	author, err := s.repo.GetAuthor()
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "查询作者信息失败")
	}

	articleCount, err := s.repo.CountArticles(author.ID)
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "统计文章数量失败")
	}

	categoryCount, err := s.repo.CountCategories(author.ID)
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "统计分类数量失败")
	}

	tagCount, err := s.repo.CountTags(author.ID)
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "统计标签数量失败")
	}

	return &dto.AuthorDTO{
		ID:            author.ID,
		Nickname:      author.Nickname,
		Bio:           author.Bio,
		Avatar:        author.Avatar,
		Github:        author.Github,
		ArticleCount:  articleCount,
		CategoryCount: categoryCount,
		TagCount:      tagCount,
	}, nil
}

// GetAbout 获取关于页面信息
func (s *authorService) GetAbout() (*dto.AboutDTO, error) {
	author, err := s.repo.GetAuthor()
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "查询作者信息失败")
	}

	return &dto.AboutDTO{
		ID:         author.ID,
		Account:    author.Account,
		Nickname:   author.Nickname,
		Avatar:     author.Avatar,
		Bio:        author.Bio,
		Github:     author.Github,
		About:      author.About,
		Background: author.Background,
		AboutBlog:  author.AboutBlog,
	}, nil
}

// AdminUpdateAuthor 后台更新作者信息
func (s *authorService) AdminUpdateAuthor(req request.AdminUpdateAuthorRequest) error {
	author, err := s.repo.GetAuthor()
	if err != nil {
		return errors.New(errors.CodeInternalError, "查询作者信息失败")
	}

	if req.Nickname != nil {
		author.Nickname = *req.Nickname
	}
	if req.Avatar != nil {
		author.Avatar = *req.Avatar
	}
	if req.Bio != nil {
		author.Bio = *req.Bio
	}
	if req.Github != nil {
		author.Github = *req.Github
	}
	if req.About != nil {
		author.About = *req.About
	}
	if req.Background != nil {
		author.Background = *req.Background
	}
	if req.AboutBlog != nil {
		author.AboutBlog = *req.AboutBlog
	}

	if err := s.repo.Update(author); err != nil {
		return errors.New(errors.CodeInternalError, "更新作者信息失败")
	}
	return nil
}

// AdminUpdateAccount 修改登录账号
func (s *authorService) AdminUpdateAccount(authorID uint, req request.AdminUpdateAccountRequest) error {
	author, err := s.repo.FindByID(authorID)
	if err != nil {
		return errors.New(errors.CodeInternalError, "查询作者信息失败")
	}

	// 检查新账号是否已存在
	if req.Account != author.Account {
		existing, _ := s.repo.FindByAccount(req.Account)
		if existing != nil {
			return errors.New(errors.CodeConflict, "该账号已存在")
		}
	}

	author.Account = req.Account
	if err := s.repo.Update(author); err != nil {
		return errors.New(errors.CodeInternalError, "修改账号失败")
	}
	return nil
}

// AdminUpdatePassword 修改登录密码
func (s *authorService) AdminUpdatePassword(authorID uint, req request.AdminUpdatePasswordRequest) error {
	author, err := s.repo.FindByID(authorID)
	if err != nil {
		return errors.New(errors.CodeInternalError, "查询作者信息失败")
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(author.Password), []byte(req.OldPassword)); err != nil {
		return errors.New(errors.CodeInvalidParam, "旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New(errors.CodeInternalError, "密码加密失败")
	}

	author.Password = string(hashedPassword)
	if err := s.repo.Update(author); err != nil {
		return errors.New(errors.CodeInternalError, "修改密码失败")
	}
	return nil
}
