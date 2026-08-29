package request

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required,min=1,max=200"`
	ParentID *uint  `json:"parent_id"` // 为空表示一级评论
}

// CommentListRequest 评论列表请求
type CommentListRequest struct {
	Page int `form:"page" binding:"required,min=1"`
	Size int `form:"size" binding:"required,min=1,max=50"`
}

// AdminCommentListRequest 后台评论列表请求
type AdminCommentListRequest struct {
	Page      int   `form:"page" binding:"required,min=1"`
	Size      int   `form:"size" binding:"required,min=1,max=50"`
	ArticleID *uint `form:"article_id"` // 按文章筛选
}
