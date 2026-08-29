package service

import (
	"blog/internal/model/dto/request"
	dto "blog/internal/model/dto/response"
	"blog/pkg/response"
)

// CommentService 评论服务接口
type CommentService interface {
	// CreateComment 创建评论
	CreateComment(articleID, userID uint, nickname, content string, parentID *uint) (*dto.CommentItem, error)
	// ListComments 获取文章的一级评论列表
	ListComments(articleID uint, page, size int, userID *uint) (*dto.CommentListResponse, error)
	// ListReplies 获取评论的回复列表
	ListReplies(commentID uint, page, size int, userID *uint) (*dto.CommentListResponse, error)
	// LikeComment 点赞/取消点赞评论
	LikeComment(commentID, userID uint) (*dto.LikeResponse, error)
	// AdminListComments 后台获取评论列表
	AdminListComments(req request.AdminCommentListRequest) (*response.PageResponse, error)
	// AdminDeleteComment 后台硬删除评论
	AdminDeleteComment(id uint) error
}
