package repository

import (
	"blog/internal/model/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// articleRepository 文章仓储实现
type articleRepository struct {
	db *gorm.DB
}

// NewArticleRepository 创建文章仓储
func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

// ListPublished 查询已发布的文章列表（按发布时间倒序）
func (r *articleRepository) ListPublished(offset, limit int) ([]entity.Article, int64, error) {
	var articles []entity.Article
	var total int64

	// 统计总数
	if err := r.db.Model(&entity.Article{}).
		Where("status = ?", 1).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询列表（按发布时间倒序）
	if err := r.db.
		Where("status = ?", 1).
		Order("published_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// SearchByKeyword 全文搜索文章（标题、描述、正文）
func (r *articleRepository) SearchByKeyword(keyword string, offset, limit int) ([]entity.Article, int64, error) {
	var articles []entity.Article
	var total int64

	query := r.db.Model(&entity.Article{}).
		Where("status = ? AND MATCH(title, summary, content) AGAINST(? IN BOOLEAN MODE)", 1, keyword)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("published_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// ListHot 查询热门文章列表（按热门权重倒序）
func (r *articleRepository) ListHot(offset, limit int) ([]entity.Article, int64, error) {
	var articles []entity.Article
	var total int64

	// 统计总数
	if err := r.db.Model(&entity.Article{}).
		Where("status = ?", 1).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询列表（按热门权重倒序）
	if err := r.db.
		Where("status = ?", 1).
		Order("hot DESC").
		Offset(offset).
		Limit(limit).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// GetArticleTagNames 查询多篇文章的标签名称（article_id -> []tag_name）
func (r *articleRepository) GetArticleTagNames(articleIDs []uint) (map[uint][]string, error) {
	if len(articleIDs) == 0 {
		return nil, nil
	}

	var results []struct {
		ArticleID uint
		TagName   string
	}
	if err := r.db.Table("article_tags").
		Select("article_tags.article_id AS article_id, tags.name AS tag_name").
		Joins("JOIN tags ON tags.id = article_tags.tag_id").
		Where("article_tags.article_id IN ?", articleIDs).
		Scan(&results).Error; err != nil {
		return nil, err
	}

	tagMap := make(map[uint][]string, len(articleIDs))
	for _, r1 := range results {
		tagMap[r1.ArticleID] = append(tagMap[r1.ArticleID], r1.TagName)
	}
	return tagMap, nil
}

// GetArticleByID 根据 ID 获取文章详情（预加载分类）
func (r *articleRepository) GetArticleByID(id uint) (*entity.Article, error) {
	var article entity.Article
	if err := r.db.Preload("Category").First(&article, id).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

// GetArticleTags 获取文章的标签列表
func (r *articleRepository) GetArticleTags(articleID uint) ([]entity.Tag, error) {
	var tags []entity.Tag
	err := r.db.Table("tags").
		Joins("JOIN article_tags ON article_tags.tag_id = tags.id").
		Where("article_tags.article_id = ?", articleID).
		Find(&tags).Error
	return tags, err
}

// HasLiked 查询用户是否已点赞某文章
func (r *articleRepository) HasLiked(db *gorm.DB, articleID, userID uint) (bool, error) {
	var count int64
	err := db.Model(&entity.ArticleLike{}).Where("article_id = ? AND user_id = ?", articleID, userID).
		Count(&count).Error
	return count > 0, err
}

// CreateLike 创建点赞记录；仅实际插入成功才同步 like_count +1（并发重复点赞不重复计数）
func (r *articleRepository) CreateLike(db *gorm.DB, articleID, userID uint) error {
	res := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&entity.ArticleLike{
		ArticleID: articleID,
		UserID:    userID,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return nil // 已点赞过，不重复计数
	}
	return db.Model(&entity.Article{}).Where("id = ?", articleID).
		UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
}

// DeleteLike 删除点赞记录；仅实际删除成功才同步 like_count -1（并发重复取消不重复扣减）
func (r *articleRepository) DeleteLike(db *gorm.DB, articleID, userID uint) error {
	res := db.Where("article_id = ? AND user_id = ?", articleID, userID).
		Delete(&entity.ArticleLike{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return nil // 未点赞过，不重复扣减
	}
	return db.Model(&entity.Article{}).Where("id = ?", articleID).
		UpdateColumn("like_count", gorm.Expr("like_count - 1")).Error
}

// GetLikeCount 获取最新点赞数
func (r *articleRepository) GetLikeCount(db *gorm.DB, articleID uint) (int, error) {
	var article entity.Article
	if err := db.Select("like_count").Where("id = ?", articleID).First(&article).Error; err != nil {
		return 0, err
	}
	return article.LikeCount, nil
}
