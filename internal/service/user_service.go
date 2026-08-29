package service

import (
	dto "blog/internal/model/dto/response"
	"blog/internal/repository"
	"blog/pkg/errors"
	"blog/pkg/response"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) AdminListUsers(page, size int) (*response.PageResponse, error) {
	offset := (page - 1) * size

	users, total, err := s.repo.AdminList(offset, size)
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "查询用户列表失败")
	}

	// 收集用户 ID
	userIDs := make([]uint, 0, len(users))
	for _, u := range users {
		userIDs = append(userIDs, u.ID)
	}

	// 批量查询评论数
	commentCounts, err := s.repo.BatchCountComments(userIDs)
	if err != nil {
		commentCounts = make(map[uint]int64)
	}

	// 批量查询点赞数
	likeCounts, err := s.repo.BatchSumLikeCount(userIDs)
	if err != nil {
		likeCounts = make(map[uint]int64)
	}

	// 组装 DTO
	list := make([]dto.AdminUserItem, 0, len(users))
	for _, u := range users {
		list = append(list, dto.AdminUserItem{
			ID:           u.ID,
			Nickname:     u.Username,
			Status:       u.Status,
			Email:        u.Email,
			RegisteredAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
			CommentCount: commentCounts[u.ID],
			LikeCount:    likeCounts[u.ID],
		})
	}

	return response.NewPageResponse(list, total, page, size), nil
}
