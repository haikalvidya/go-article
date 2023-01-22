package usecase

import (
	"errors"

	"github.com/haikalvidya/go-article/internal/delivery/payload"
	"github.com/haikalvidya/go-article/internal/models"
	"gorm.io/gorm"
)

type IArticleUsecase interface {
	CreateArticle(authorID string, req *payload.CreateArticleRequest) (*payload.ArticleInfo, error)
	GetAllArticles() ([]*payload.ArticleInfo, error)
	GetArticleByID(id int) (*payload.ArticleInfo, error)
	GetArticlesByAuthorID(authorID string) ([]*payload.ArticleInfo, error)
	SearchArticlesByTitleAndContent(content string) ([]*payload.ArticleInfo, error)
	DeleteArticleByID(id int, authorId string) error
	UpdateArticleByID(id int, req *payload.CreateArticleRequest, authorId string) (*payload.ArticleInfo, error)
}

type articleUsecase usecaseType

func (u *articleUsecase) CreateArticle(authorID string, req *payload.CreateArticleRequest) (*payload.ArticleInfo, error) {
	// using createtx
	article := &models.ArticleModel{
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: authorID,
	}

	err := u.Repo.Tx.DoInTransaction(func(tx *gorm.DB) error {
		createdArticle, err := u.Repo.Article.CreateTx(tx, article)
		if err != nil {
			return err
		}

		article = createdArticle
		return nil
	})
	if err != nil {
		return nil, err
	}

	return article.PublicInfo(), nil
}

func (u *articleUsecase) GetAllArticles() ([]*payload.ArticleInfo, error) {
	articles, err := u.Repo.Article.GetAll()
	if err != nil {
		return nil, err
	}

	res := make([]*payload.ArticleInfo, 0)
	for _, article := range articles {
		res = append(res, article.PublicInfo())
	}

	return res, nil
}

func (u *articleUsecase) GetArticleByID(id int) (*payload.ArticleInfo, error) {
	article, err := u.Repo.Article.SelectByID(id)
	if err != nil {
		return nil, err
	}

	return article.PublicInfo(), nil
}

func (u *articleUsecase) GetArticlesByAuthorID(authorID string) ([]*payload.ArticleInfo, error) {
	articles, err := u.Repo.Article.SelectByAuthorID(authorID)
	if err != nil {
		return nil, err
	}

	res := make([]*payload.ArticleInfo, 0)
	for _, article := range articles {
		res = append(res, article.PublicInfo())
	}

	return res, nil
}

func (u *articleUsecase) SearchArticlesByTitleAndContent(content string) ([]*payload.ArticleInfo, error) {
	articles, err := u.Repo.Article.SearchByTitleAndContent(content)
	if err != nil {
		return nil, err
	}

	res := make([]*payload.ArticleInfo, 0)
	for _, article := range articles {
		res = append(res, article.PublicInfo())
	}

	return res, nil
}

func (u *articleUsecase) DeleteArticleByID(id int, authorId string) error {
	// check if article is owned by author
	article, err := u.Repo.Article.SelectByID(id)
	if err != nil {
		return err
	}

	if article.AuthorID != authorId {
		return errors.New(payload.ERROR_ARTICLE_NOT_ALLOWED)
	}

	err = u.Repo.Tx.DoInTransaction(func(tx *gorm.DB) error {
		article := &models.ArticleModel{
			ID: id,
		}
		err := u.Repo.Article.DeleteTx(tx, article)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

func (u *articleUsecase) UpdateArticleByID(id int, req *payload.CreateArticleRequest, authorId string) (*payload.ArticleInfo, error) {
	// check if article is owned by author
	article, err := u.Repo.Article.SelectByID(id)
	if err != nil {
		return nil, err
	}

	if article.AuthorID != authorId {
		return nil, errors.New(payload.ERROR_ARTICLE_NOT_ALLOWED)
	}

	err = u.Repo.Tx.DoInTransaction(func(tx *gorm.DB) error {
		err := u.Repo.Article.UpdateTx(tx, article)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return article.PublicInfo(), nil
}
