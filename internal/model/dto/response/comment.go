package dto

import "time"

// CommentItem 评论项
type CommentItem struct {
	ID         uint      `json:"id"`
	Nickname   string    `json:"nickname"`
	Content    string    `json:"content"`
	LikeCount  int       `json:"like_count"`
	Liked      bool      `json:"liked"`
	CreatedAt  time.Time `json:"created_at"`
	ReplyCount int       `json:"reply_count"` // 子评论数（一级评论用）
}

// CommentListResponse 评论列表响应
type CommentListResponse struct {
	List  []CommentItem `json:"list"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Size  int           `json:"size"`
}
