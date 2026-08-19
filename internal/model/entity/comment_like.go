package entity

// CommentLike 评论点赞记录（按用户去重，支持取消点赞与「已点赞」状态）
// 与文章点赞 article_likes 同构，评论的点赞数由 Comment.LikeCount 冗余存储。
type CommentLike struct {
	Base
	CommentID uint `gorm:"column:comment_id;type:bigint;not null;uniqueIndex:uk_comment_user" json:"comment_id"`
	UserID    uint `gorm:"column:user_id;type:bigint;not null;uniqueIndex:uk_comment_user" json:"user_id"`
}

func (CommentLike) TableName() string { return "comment_likes" }
