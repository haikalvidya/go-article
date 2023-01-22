package repository

import (
	"github.com/haikalvidya/go-article/internal/models"
	"gorm.io/gorm"
)

type IArticleRepository interface {
	GetAll() ([]*models.ArticleModel, error)
	SelectByID(id int) (*models.ArticleModel, error)
	SelectByAuthorID(authorID string) ([]*models.ArticleModel, error)
	SearchByTitleAndContent(content string) ([]*models.ArticleModel, error)
	CreateTx(tx *gorm.DB, article *models.ArticleModel) (*models.ArticleModel, error)
	DeleteTx(tx *gorm.DB, article *models.ArticleModel) error
	UpdateTx(tx *gorm.DB, article *models.ArticleModel) error
}

type articleRepository repositoryType

func (r *articleRepository) GetAll() ([]*models.ArticleModel, error) {
	articles := []*models.ArticleModel{}
	err := r.DB.Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (r *articleRepository) SelectByID(id int) (*models.ArticleModel, error) {
	article := &models.ArticleModel{}
	err := r.DB.Where("id = ?", id).First(article).Error
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (r *articleRepository) SelectByAuthorID(authorID string) ([]*models.ArticleModel, error) {
	articles := []*models.ArticleModel{}
	err := r.DB.Where("author_id = ?", authorID).Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (r *articleRepository) SearchByTitleAndContent(content string) ([]*models.ArticleModel, error) {
	articles := []*models.ArticleModel{}
	err := r.DB.Where("title LIKE ? OR content LIKE ?", "%"+content+"%", "%"+content+"%").Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (r *articleRepository) CreateTx(tx *gorm.DB, article *models.ArticleModel) (*models.ArticleModel, error) {
	err := tx.Create(article).Error
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (r *articleRepository) DeleteTx(tx *gorm.DB, article *models.ArticleModel) error {
	err := tx.Delete(&article).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *articleRepository) UpdateTx(tx *gorm.DB, article *models.ArticleModel) error {
	err := tx.Save(&article).Error
	if err != nil {
		return err
	}
	return nil
}
