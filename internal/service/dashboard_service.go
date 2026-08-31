package service

import (
	dto "blog/internal/model/dto/response"
	"blog/internal/repository"
	"blog/pkg/errors"
)

type dashboardService struct {
	articleRepo repository.ArticleRepository
	commentRepo repository.CommentRepository
	userRepo    repository.UserRepository
}

func NewDashboardService(
	articleRepo repository.ArticleRepository,
	commentRepo repository.CommentRepository,
	userRepo repository.UserRepository,
) DashboardService {
	return &dashboardService{
		articleRepo: articleRepo,
		commentRepo: commentRepo,
		userRepo:    userRepo,
	}
}

func (s *dashboardService) GetDashboard() (*dto.DashboardDTO, error) {
	// 统计文章总数
	articleCount, err := s.articleRepo.Count()
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "统计文章数量失败")
	}

	// 统计评论总数
	commentCount, err := s.commentRepo.Count()
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "统计评论数量失败")
	}

	// 统计用户总数
	userCount, err := s.userRepo.Count()
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "统计用户数量失败")
	}

	return &dto.DashboardDTO{
		ArticleCount: articleCount,
		CommentCount: commentCount,
		UserCount:    userCount,
	}, nil
}
