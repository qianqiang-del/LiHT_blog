package entity

import "time"

// Article 文章
type Article struct {
	Base
	Title        string     `gorm:"column:title;type:varchar(200);not null;index" json:"title"`
	Summary      string     `gorm:"column:summary;type:varchar(500)" json:"summary"`
	Content      string     `gorm:"column:content;type:longtext" json:"content"`
	Cover        string     `gorm:"column:cover;type:varchar(255)" json:"cover"`
	AuthorID     uint       `gorm:"column:author_id;type:bigint;index" json:"author_id"`
	CategoryID   uint       `gorm:"column:category_id;type:bigint;index" json:"category_id"`
	Status       int8       `gorm:"column:status;type:tinyint;not null;default:1" json:"status"` // 1=已发布 0=下架(隐藏)
	ViewCount    int        `gorm:"column:view_count;type:int;not null;default:0" json:"view_count"`
	LikeCount    int        `gorm:"column:like_count;type:int;not null;default:0" json:"like_count"`
	CommentCount int        `gorm:"column:comment_count;type:int;not null;default:0" json:"comment_count"`
	Hot          int        `gorm:"column:hot;type:int;not null;default:0;index" json:"hot"` // 热门权重：越大越热，按照规则计算
	PublishedAt  *time.Time `gorm:"column:published_at" json:"published_at"`
}

func (Article) TableName() string { return "articles" }
