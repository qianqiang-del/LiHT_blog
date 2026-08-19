package entity

// User 用户/会员（前台登录、注册、评论作者、文章作者共用）
type User struct {
	Base
	Username string `gorm:"column:username;type:varchar(50);not null;uniqueIndex" json:"username"`
	Email    string `gorm:"column:email;type:varchar(120);uniqueIndex" json:"email"`
	Password string `gorm:"column:password;type:varchar(255);not null" json:"-"` // 存储哈希，不对外暴露
	Nickname string `gorm:"column:nickname;type:varchar(50)" json:"nickname"`
	Avatar   string `gorm:"column:avatar;type:varchar(255)" json:"avatar"`
	Bio      string `gorm:"column:bio;type:varchar(500)" json:"bio"`
	Github   string `gorm:"column:github;type:varchar(255)" json:"github"`
	Role     int8   `gorm:"column:role;type:tinyint;not null;default:1" json:"role"`     // 1=普通用户 2=管理员
	Status   int8   `gorm:"column:status;type:tinyint;not null;default:1" json:"status"` // 1=正常 0=禁用
}

func (User) TableName() string { return "users" }
