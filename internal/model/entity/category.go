package entity

// Category 文章分类
type Category struct {
	Base
	Name string `gorm:"column:name;type:varchar(50);not null;uniqueIndex" json:"name"`
}

func (Category) TableName() string { return "categories" }
