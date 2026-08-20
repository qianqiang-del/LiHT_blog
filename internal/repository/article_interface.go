package repository

import "blog/internal/model/entity"

// ArticleRepository 文章仓储接口
type ArticleRepository interface {
	ListPublished(offset, limit int) ([]entity.Article, int64, error)
	ListHot(offset, limit int) ([]entity.Article, int64, error)
	// GetArticleTagNames 查询多篇文章的标签名称（article_id -> []tag_name）
	GetArticleTagNames(articleIDs []uint) (map[uint][]string, error)
}
