package dto

// DashboardDTO 概览数据
type DashboardDTO struct {
	ArticleCount int64 `json:"article_count"`
	CommentCount int64 `json:"comment_count"`
	UserCount    int64 `json:"user_count"`
}
