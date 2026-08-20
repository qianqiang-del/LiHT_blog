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
