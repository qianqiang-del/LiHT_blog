package dto

import "time"

// ArticleListItem 文章列表项
type ArticleListItem struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Summary     string     `json:"summary"`
	Cover       string     `json:"cover"`
	Tags        []string   `json:"tags"`
	ViewCount   int        `json:"view_count"`
	PublishedAt *time.Time `json:"published_at"`
}

// HotArticleItem 热门文章列表项
type HotArticleItem struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Summary     string     `json:"summary"`
	Cover       string     `json:"cover"`
	PublishedAt *time.Time `json:"published_at"`
}

// ArticleDetail 文章详情
type ArticleDetail struct {
	ID           uint         `json:"id"`
	Title        string       `json:"title"`
	Summary      string       `json:"summary"`
	Content      string       `json:"content"`
	Cover        string       `json:"cover"`
	Category     *CategoryDTO `json:"category,omitempty"`
	Tags         []TagDTO     `json:"tags"`
	ViewCount    int          `json:"view_count"`
	LikeCount    int          `json:"like_count"`
	CommentCount int          `json:"comment_count"`
	PublishedAt  *time.Time   `json:"published_at"`
	Liked        bool         `json:"liked"`
}

// LikeResponse 点赞操作响应
type LikeResponse struct {
	Liked     bool `json:"liked"`
	LikeCount int  `json:"like_count"`
}

// AdminArticleListItem 后台文章列表项
type AdminArticleListItem struct {
	ID           uint       `json:"id"`
	Title        string     `json:"title"`
	Cover        string     `json:"cover"`
	Category     string     `json:"category"`
	Tags         []string   `json:"tags"`
	Status       int8       `json:"status"` // 1=正常 0=下架
	ViewCount    int        `json:"view_count"`
	LikeCount    int        `json:"like_count"`
	CommentCount int        `json:"comment_count"`
	Hot          int        `json:"hot"`
	PublishedAt  *time.Time `json:"published_at"`
}

// ArticleOption 文章选项（用于下拉选择）
type ArticleOption struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}
