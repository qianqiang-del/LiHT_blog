package entity

import "gorm.io/gorm"

// Comment 评论（支持嵌套回复，parent_id 指向父评论）
type Comment struct {
	Base
	ArticleID uint           `gorm:"column:article_id;type:bigint;index" json:"article_id"`
	ParentID  *uint          `gorm:"column:parent_id;type:bigint" json:"parent_id"` // 为空表示一级评论
	UserID    *uint          `gorm:"column:user_id;type:bigint" json:"user_id"`     // 登录用户ID，游客为空
	Nickname  string         `gorm:"column:nickname;type:varchar(50);not null" json:"nickname"`
	Content   string         `gorm:"column:content;type:varchar(1000);not null" json:"content"`
	LikeCount int            `gorm:"column:like_count;type:int;not null;default:0" json:"like_count"`
	Status    int8           `gorm:"column:status;type:tinyint;not null;default:1" json:"status"` // 1=正常 0=删除
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (Comment) TableName() string { return "comments" }
