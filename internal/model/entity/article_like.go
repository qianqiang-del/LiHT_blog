package entity

// ArticleLike 文章点赞记录（按用户去重，支持取消点赞与「已点赞」状态）
type ArticleLike struct {
	Base
	ArticleID uint `gorm:"column:article_id;type:bigint;not null;uniqueIndex:uk_article_user" json:"article_id"`
	UserID    uint `gorm:"column:user_id;type:bigint;not null;uniqueIndex:uk_article_user" json:"user_id"`
}

func (ArticleLike) TableName() string { return "article_likes" }
