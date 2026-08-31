package entity

// Author 作者/博主
type Author struct {
	Base
	Account    string `gorm:"column:account;type:varchar(50);not null;uniqueIndex" json:"account"` // 登录账号
	Password   string `gorm:"column:password;type:varchar(255);not null" json:"-"`                 // 登录密码（哈希存储）
	Nickname   string `gorm:"column:nickname;type:varchar(50);not null" json:"nickname"`
	Avatar     string `gorm:"column:avatar;type:varchar(255)" json:"avatar"`
	Bio        string `gorm:"column:bio;type:varchar(500)" json:"bio"`
	Github     string `gorm:"column:github;type:varchar(255)" json:"github"`
	About      string `gorm:"column:about;type:text" json:"about"`                   // 关于我
	Background string `gorm:"column:background;type:varchar(255)" json:"background"` // 背景图
	AboutBlog  string `gorm:"column:about_blog;type:text" json:"about_blog"`         // 关于博客
}

func (Author) TableName() string { return "authors" }
