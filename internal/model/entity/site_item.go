package entity

// SiteItem 关于页可配置内容统一存储（技术栈分组 / 技术栈项 / 联系方式 / 博客故事）
// type 取值：skill_group | skill | contact | story
//   - skill_group：title=分组名，其余为空
//   - skill：title=技能名，parent_id=所属分组
//   - contact：title=展示标签，content=值，code=类型(github/email/wechat)，icon=图标
//   - story：title=时间文本，content=故事内容
type SiteItem struct {
	Base
	Type     string `gorm:"column:type;type:varchar(30);not null;index" json:"type"`
	Title    string `gorm:"column:title;type:varchar(200);not null" json:"title"`
	Content  string `gorm:"column:content;type:varchar(500)" json:"content"`
	ParentID *uint  `gorm:"column:parent_id;type:bigint" json:"parent_id"` // skill -> 所属分组
	Code     string `gorm:"column:code;type:varchar(30)" json:"code"`      // 联系类型
	Icon     string `gorm:"column:icon;type:varchar(50)" json:"icon"`      // 联系图标
	Sort     int    `gorm:"column:sort;type:int;not null;default:0" json:"sort"`
}

func (SiteItem) TableName() string { return "site_items" }
