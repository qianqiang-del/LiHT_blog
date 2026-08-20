package dto

// AuthorDTO 作者详情
type AuthorDTO struct {
	ID            uint   `json:"id"`
	Nickname      string `json:"nickname"`
	Bio           string `json:"bio"`
	Avatar        string `json:"avatar"`
	Github        string `json:"github"`
	ArticleCount  int64  `json:"article_count"`
	CategoryCount int64  `json:"category_count"`
	TagCount      int64  `json:"tag_count"`
}
