package entity

// Tag 文章标签
type Tag struct {
	Base
	Name string `gorm:"column:name;type:varchar(50);not null;uniqueIndex" json:"name"`
}

func (Tag) TableName() string { return "tags" }
