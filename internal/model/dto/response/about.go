package dto

// AboutDTO 关于页面
type AboutDTO struct {
	ID         uint   `json:"id"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	Bio        string `json:"bio"`
	Github     string `json:"github"`
	About      string `json:"about"`
	Background string `json:"background"`
	AboutBlog  string `json:"about_blog"`
}
