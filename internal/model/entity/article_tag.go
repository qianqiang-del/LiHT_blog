package entity

// ArticleTag 文章-标签 多对多关联表
type ArticleTag struct {
	ArticleID uint `gorm:"column:article_id;primaryKey" json:"article_id"`
	TagID     uint `gorm:"column:tag_id;primaryKey" json:"tag_id"`
}

func (ArticleTag) TableName() string { return "article_tags" }
