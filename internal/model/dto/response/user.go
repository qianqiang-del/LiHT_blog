package dto

// AdminUserItem 后台用户列表项
type AdminUserItem struct {
	ID           uint   `json:"id"`
	Nickname     string `json:"nickname"`
	Status       int8   `json:"status"`
	Email        string `json:"email"`
	RegisteredAt string `json:"registered_at"`
	CommentCount int64  `json:"comment_count"`
	LikeCount    int64  `json:"like_count"`
}
