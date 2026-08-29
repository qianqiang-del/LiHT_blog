package service

import "blog/pkg/response"

// UserService 用户服务接口
type UserService interface {
	// AdminListUsers 后台获取用户列表
	AdminListUsers(page, size int) (*response.PageResponse, error)
}
