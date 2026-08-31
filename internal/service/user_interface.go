package service

import (
	"blog/internal/model/dto/request"
	"blog/pkg/response"
)

// UserService 用户服务接口
type UserService interface {
	// AdminListUsers 后台获取用户列表
	AdminListUsers(page, size int) (*response.PageResponse, error)
	// AdminUpdateUserStatus 后台更新用户状态
	AdminUpdateUserStatus(id uint, req request.AdminUpdateUserStatusRequest) error
	// AdminDeleteUser 后台删除用户
	AdminDeleteUser(id uint) error
}
