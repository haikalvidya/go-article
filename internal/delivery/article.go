package delivery

import (
	"net/http"
	"strconv"

	"github.com/haikalvidya/go-article/internal/delivery/payload"
	"github.com/haikalvidya/go-article/pkg/common"
	"github.com/haikalvidya/go-article/pkg/utils"
	"github.com/labstack/echo/v4"
)

type articleDelivery deliveryType

func (d *articleDelivery) CreateArticle(c echo.Context) error {
	res := common.Response{}
	req := &payload.CreateArticleRequest{}

	userId := d.Middleware.JWT.GetUserIdFromJwt(c)

	c.Bind(req)

	if err := c.Validate(req); err != nil {
		res.Error = utils.GetErrorValidation(err)
		res.Status = false
		res.Message = "Failed Create Article"
		return c.JSON(http.StatusBadRequest, res)
	}

	articleRes, err := d.Usecase.Article.CreateArticle(userId, req)
	if err != nil {
		res.Status = false
		res.Message = err.Error()
		return c.JSON(http.StatusBadRequest, res)
	}

	res.Message = "Success Create Article"
	res.Data = articleRes
	res.Status = true

	return c.JSON(http.StatusOK, res)
}

// get all article
func (d *articleDelivery) GetAllArticle(c echo.Context) error {
	res := common.Response{}
	queryParam := &payload.ArticleQuery{}

	var err error

	// bind query param
	if err := c.Bind(queryParam); err != nil {
		res.Status = false
		res.Message = err.Error()
		return c.JSON(http.StatusBadRequest, res)
	}

	articleRes := []*payload.ArticleInfo{}

	if queryParam.AuthorName != "" && queryParam.QuerySearch != "" {
		userRes, err := d.Usecase.User.GetUserByName(queryParam.AuthorName)
		if err != nil {
			res.Status = false
			res.Message = err.Error()
			res.Data = []*payload.ArticleInfo{}
			return c.JSON(http.StatusBadRequest, res)
		}
		articleRes, err = d.Usecase.Article.GetArticleSearchAndByAuthorID(userRes.ID, queryParam.QuerySearch)
		if err != nil {
			res.Status = false
			res.Message = err.Error()
			res.Data = []*payload.ArticleInfo{}
			return c.JSON(http.StatusBadRequest, res)
		}
	} else if queryParam.QuerySearch != "" {
		articleRes, err = d.Usecase.Article.SearchArticlesByTitleAndContent(queryParam.QuerySearch)
		if err != nil {
			res.Status = false
			res.Message = err.Error()
			res.Data = []*payload.ArticleInfo{}
			return c.JSON(http.StatusBadRequest, res)
		}
	} else if queryParam.AuthorName != "" {
		// get user id by author name
		userRes, err := d.Usecase.User.GetUserByName(queryParam.AuthorName)
		if err != nil {
			res.Status = false
			res.Message = err.Error()
			res.Data = []*payload.ArticleInfo{}
			return c.JSON(http.StatusBadRequest, res)
		}

		articleRes, err = d.Usecase.Article.GetArticlesByAuthorID(userRes.ID)
		if err != nil {
			res.Status = false
			res.Message = err.Error()
			res.Data = []*payload.ArticleInfo{}
			return c.JSON(http.StatusBadRequest, res)
		}
	} else {
		articleRes, err = d.Usecase.Article.GetAllArticles()
		if err != nil {
			res.Status = false
			res.Message = err.Error()
			res.Data = []*payload.ArticleInfo{}
			return c.JSON(http.StatusBadRequest, res)
		}
	}
	res.Message = "Success Get All Article"
	res.Data = articleRes
	res.Status = true

	return c.JSON(http.StatusOK, res)
}

// get article by id
func (d *articleDelivery) GetArticleByID(c echo.Context) error {
	res := common.Response{}
	articleIDStr := c.Param("id")

	// convert string to int
	articleID, _ := strconv.Atoi(articleIDStr)

	articleRes, err := d.Usecase.Article.GetArticleByID(articleID)
	if err != nil {
		res.Status = false
		res.Message = err.Error()
		return c.JSON(http.StatusBadRequest, res)
	}

	res.Message = "Success Get Article By ID"
	res.Data = articleRes
	res.Status = true

	return c.JSON(http.StatusOK, res)
}

// update article
func (d *articleDelivery) UpdateArticle(c echo.Context) error {
	res := common.Response{}
	req := &payload.UpdateArticleRequest{}
	articleIDStr := c.Param("id")

	// convert string to int
	articleID, _ := strconv.Atoi(articleIDStr)

	userId := d.Middleware.JWT.GetUserIdFromJwt(c)

	c.Bind(req)

	if err := c.Validate(req); err != nil {
		res.Error = utils.GetErrorValidation(err)
		res.Status = false
		res.Message = "Failed Update Article"
		return c.JSON(http.StatusBadRequest, res)
	}

	articleRes, err := d.Usecase.Article.UpdateArticleByID(articleID, req, userId)
	if err != nil {
		res.Status = false
		res.Message = err.Error()
		return c.JSON(http.StatusBadRequest, res)
	}

	res.Message = "Success Update Article"
	res.Data = articleRes
	res.Status = true

	return c.JSON(http.StatusOK, res)
}

// delete article
func (d *articleDelivery) DeleteArticle(c echo.Context) error {
	res := common.Response{}
	articleIDStr := c.Param("id")

	// convert string to int
	articleID, _ := strconv.Atoi(articleIDStr)

	userId := d.Middleware.JWT.GetUserIdFromJwt(c)

	err := d.Usecase.Article.DeleteArticleByID(articleID, userId)
	if err != nil {
		res.Status = false
		res.Message = err.Error()
		return c.JSON(http.StatusBadRequest, res)
	}

	res.Message = "Success Delete Article"
	res.Status = true

	return c.JSON(http.StatusOK, res)
}
