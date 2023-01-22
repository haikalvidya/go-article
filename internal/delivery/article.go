package delivery

import (
	"net/http"

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

	articleRes, err := d.Usecase.Article.GetAllArticles()
	if err != nil {
		res.Status = false
		res.Message = err.Error()
		return c.JSON(http.StatusBadRequest, res)
	}

	res.Message = "Success Get All Article"
	res.Data = articleRes
	res.Status = true

	return c.JSON(http.StatusOK, res)
}
