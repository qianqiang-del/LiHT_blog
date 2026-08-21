package service

import dto "blog/internal/model/dto/response"

// TagService 标签服务接口
type TagService interface {
	ListTags() ([]dto.TagDTO, error)
}
