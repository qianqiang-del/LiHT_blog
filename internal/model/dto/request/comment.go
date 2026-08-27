package request

import "blog/pkg/response"

// CommentListRequest 评论列表请求
type CommentListRequest struct {
	response.PageRequest
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required,min=1,max=1000"`
	ParentID *uint  `json:"parent_id"` // 为空表示一级评论，有值表示回复
}
