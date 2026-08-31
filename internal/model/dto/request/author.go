package request

// AdminUpdateAuthorRequest 后台更新作者信息请求
type AdminUpdateAuthorRequest struct {
	Nickname   *string `json:"nickname"`
	Avatar     *string `json:"avatar"`
	Bio        *string `json:"bio"`
	Github     *string `json:"github"`
	About      *string `json:"about"`
	Background *string `json:"background"`
	AboutBlog  *string `json:"about_blog"`
}
