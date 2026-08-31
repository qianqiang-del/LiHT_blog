package service

import (
	"blog/internal/model/dto/request"
	dto "blog/internal/model/dto/response"
)

// AuthorService 作者服务接口
type AuthorService interface {
	GetAuthorDetail() (*dto.AuthorDTO, error)
	GetAbout() (*dto.AboutDTO, error)
	AdminUpdateAuthor(req request.AdminUpdateAuthorRequest) error
}
