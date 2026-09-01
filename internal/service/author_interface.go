package service

import (
	"blog/internal/model/dto/request"
	dto "blog/internal/model/dto/response"
)

// AuthorService 作者服务接口
type AuthorService interface {
	GetAuthorDetail() (*dto.AuthorDTO, error)
	GetAbout() (*dto.AboutDTO, error)
	AdminLogin(req request.AdminLoginRequest) (*dto.AdminLoginDTO, error)
	AdminLogout(token string) error
	AdminGetCaptcha() (*dto.CaptchaResponse, error)
	AdminUpdateAuthor(req request.AdminUpdateAuthorRequest) error
	AdminUpdateAccount(authorID uint, req request.AdminUpdateAccountRequest) error
	AdminUpdatePassword(authorID uint, req request.AdminUpdatePasswordRequest) error
}
