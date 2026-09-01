package repository

import (
	"blog/internal/model/entity"
	"gorm.io/gorm"
)

// authorRepository 作者详情仓储实现
type authorRepository struct {
	db *gorm.DB
}

// NewAuthorRepository 创建作者仓储
func NewAuthorRepository(db *gorm.DB) AuthorRepository {
	return &authorRepository{db: db}
}

// GetAuthor 获取作者详情信息
func (r *authorRepository) GetAuthor() (*entity.Author, error) {
	var author entity.Author
	if err := r.db.First(&author).Error; err != nil {
		return nil, err
	}
	return &author, nil
}

// FindByID 根据 ID 查询作者
func (r *authorRepository) FindByID(id uint) (*entity.Author, error) {
	var author entity.Author
	if err := r.db.First(&author, id).Error; err != nil {
		return nil, err
	}
	return &author, nil
}

// FindByAccount 根据账号查询作者
func (r *authorRepository) FindByAccount(account string) (*entity.Author, error) {
	var author entity.Author
	if err := r.db.Where("account = ?", account).First(&author).Error; err != nil {
		return nil, err
	}
	return &author, nil
}

// Update 更新作者信息
func (r *authorRepository) Update(author *entity.Author) error {
	return r.db.Save(author).Error
}

// CountArticles 统计作者发布文章数量
func (r *authorRepository) CountArticles(authorID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Article{}).
		Where("author_id = ?", authorID).
		Count(&count).Error
	return count, err
}

// CountCategories 统计作者的分类数量（去重）
func (r *authorRepository) CountCategories(authorID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Article{}).
		Where("author_id = ?", authorID).
		Distinct("category_id").
		Count(&count).Error
	return count, err
}

// CountTags 统计作者的标签数量（去重）
func (r *authorRepository) CountTags(authorID uint) (int64, error) {
	var count int64
	err := r.db.Table("article_tags").
		Joins("JOIN articles ON articles.id = article_tags.article_id").
		Where("articles.author_id = ?", authorID).
		Distinct("article_tags.tag_id").
		Count(&count).Error
	return count, err
}
