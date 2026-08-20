package service

import dto "blog/internal/model/dto/response"

// AuthorService 作者服务接口
type AuthorService interface {
	GetAuthorDetail() (*dto.AuthorDTO, error)
}
