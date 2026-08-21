package dto

// TagDTO 标签
type TagDTO struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	ArticleCount int64  `json:"article_count"`
}
