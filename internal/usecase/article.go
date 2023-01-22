package usecase

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/go-redis/redis"
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
	GetArticleSearchAndByAuthorID(authorID string, content string) ([]*payload.ArticleInfo, error)
	DeleteArticleByID(id int, authorId string) error
	UpdateArticleByID(id int, req *payload.UpdateArticleRequest, authorId string) (*payload.ArticleInfo, error)
}

type articleUsecase usecaseType

func (u *articleUsecase) CreateArticle(authorID string, req *payload.CreateArticleRequest) (*payload.ArticleInfo, error) {
	// check if author exist and login
	author, err := u.Repo.User.SelectByID(authorID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if author == nil {
		return nil, errors.New(payload.ERROR_USER_NOT_FOUND)
	}

	// check in redis if user is logged in
	_, err = u.RedisClient.Get(authorID).Result()
	if err != nil {
		return nil, errors.New(payload.ERROR_USER_NOT_LOGGED_IN)
	}

	// using createtx
	article := &models.ArticleModel{
		Title:    req.Title,
		Body:     req.Content,
		AuthorID: authorID,
	}

	u.RedisClient.Del("GET_ALL_ARTICLES")

	err = u.Repo.Tx.DoInTransaction(func(tx *gorm.DB) error {
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
	data, err := u.RedisClient.Get("GET_ALL_ARTICLES").Result()
	if err != nil && err != redis.Nil {
		return nil, errors.New(payload.ERROR_GET_ARTICLE)
	}

	var res []*payload.ArticleInfo

	if data != "" {
		err := json.Unmarshal([]byte(data), &res)
		if err != nil {
			return nil, errors.New(payload.ERROR_GET_ARTICLE)
		}

		return res, nil
	} else {
		articles, err := u.Repo.Article.GetAll()
		if err != nil {
			return nil, err
		}

		res = make([]*payload.ArticleInfo, 0)
		for _, article := range articles {
			res = append(res, article.PublicInfo())
		}

		dataJsonByte, err := json.Marshal(res)
		if err != nil {
			return nil, errors.New(payload.ERROR_GET_ARTICLE)
		}
		dataJson := string(dataJsonByte)

		u.RedisClient.Set("GET_ALL_ARTICLES", dataJson, 0)
	}

	return res, nil
}

func (u *articleUsecase) GetArticleByID(id int) (*payload.ArticleInfo, error) {
	data, err := u.RedisClient.Get("GET_ARTICLE_BY_ID_" + strconv.Itoa(id)).Result()
	if err != nil && err != redis.Nil {
		return nil, errors.New(payload.ERROR_GET_ARTICLE)
	}

	var res *payload.ArticleInfo

	if data != "" {
		err := json.Unmarshal([]byte(data), &res)
		if err != nil {
			return nil, errors.New(payload.ERROR_GET_ARTICLE)
		}
	} else {
		article, err := u.Repo.Article.SelectByID(id)
		if err != nil {
			return nil, err
		}

		res = article.PublicInfo()

		dataJsonByte, err := json.Marshal(res)
		if err != nil {
			return nil, errors.New(payload.ERROR_GET_ARTICLE)
		}
		dataJson := string(dataJsonByte)
		u.RedisClient.Set("GET_ARTICLE_BY_ID_"+strconv.Itoa(id), dataJson, 0)
	}

	return res, nil
}

func (u *articleUsecase) GetArticlesByAuthorID(authorID string) ([]*payload.ArticleInfo, error) {
	data, err := u.RedisClient.Get("GET_ARTICLES_BY_AUTHOR_ID_" + authorID).Result()
	if err != nil && err != redis.Nil {
		return nil, errors.New(payload.ERROR_GET_ARTICLE)
	}

	var res []*payload.ArticleInfo

	if data != "" {
		err := json.Unmarshal([]byte(data), &res)
		if err != nil {
			return nil, errors.New(payload.ERROR_GET_ARTICLE)
		}
	} else {
		articles, err := u.Repo.Article.SelectByAuthorID(authorID)
		if err != nil {
			return nil, err
		}

		res = make([]*payload.ArticleInfo, 0)
		for _, article := range articles {
			res = append(res, article.PublicInfo())
		}

		dataJsonByte, err := json.Marshal(res)
		if err != nil {
			return nil, errors.New(payload.ERROR_GET_ARTICLE)
		}
		dataJson := string(dataJsonByte)
		u.RedisClient.Set("GET_ARTICLES_BY_AUTHOR_ID_"+authorID, dataJson, 0)
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

	// check in redis if user is logged in
	_, err = u.RedisClient.Get(authorId).Result()
	if err != nil {
		return errors.New(payload.ERROR_USER_NOT_LOGGED_IN)
	}

	// delete all redis cache about article
	u.RedisClient.Del("GET_ARTICLE_BY_ID_" + strconv.Itoa(id))
	u.RedisClient.Del("GET_ARTICLES_BY_AUTHOR_ID_" + authorId)
	u.RedisClient.Del("GET_ALL_ARTICLES")

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

func (u *articleUsecase) UpdateArticleByID(id int, req *payload.UpdateArticleRequest, authorId string) (*payload.ArticleInfo, error) {
	// check if author exist and login
	author, err := u.Repo.User.SelectByID(authorId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if author == nil {
		return nil, errors.New(payload.ERROR_USER_NOT_FOUND)
	}

	// check in redis if user is logged in
	_, err = u.RedisClient.Get(authorId).Result()
	if err != nil {
		return nil, errors.New(payload.ERROR_USER_NOT_LOGGED_IN)
	}

	// check if article is owned by author
	article, err := u.Repo.Article.SelectByID(id)
	if err != nil {
		return nil, err
	}

	article.Author = nil

	if article.AuthorID != authorId {
		return nil, errors.New(payload.ERROR_ARTICLE_NOT_ALLOWED)
	}

	// delete all redis cache about article
	u.RedisClient.Del("GET_ARTICLE_BY_ID_" + strconv.Itoa(id))
	u.RedisClient.Del("GET_ARTICLES_BY_AUTHOR_ID_" + authorId)
	u.RedisClient.Del("GET_ALL_ARTICLES")

	if req.Title != "" {
		article.Title = req.Title
	}

	if req.Content != "" {
		article.Body = req.Content
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

	article.Author = author

	return article.PublicInfo(), nil
}

func (u *articleUsecase) GetArticleSearchAndByAuthorID(authorID string, content string) ([]*payload.ArticleInfo, error) {
	articles, err := u.Repo.Article.SearchByTitleAndContentAndAuthorID(authorID, content)
	if err != nil {
		return nil, err
	}

	res := make([]*payload.ArticleInfo, 0)
	for _, article := range articles {
		res = append(res, article.PublicInfo())
	}

	return res, nil
}
