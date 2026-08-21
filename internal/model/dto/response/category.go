package dto

// CategoryDTO 分类
type CategoryDTO struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	ArticleCount int64  `json:"article_count"`
}
