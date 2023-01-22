package delivery

import (
	"github.com/haikalvidya/go-article/internal/middlewares"
	"github.com/haikalvidya/go-article/internal/usecase"

	"github.com/labstack/echo/v4"
)

type Delivery struct {
	User    *userDelivery
	Article *articleDelivery
}

type deliveryType struct {
	Usecase    *usecase.Usecase
	Middleware *middlewares.CustomMiddleware
}

func NewDelivery(e *echo.Echo, usecase *usecase.Usecase, mid *middlewares.CustomMiddleware) *Delivery {
	deliveryType := &deliveryType{
		Usecase:    usecase,
		Middleware: mid,
	}
	delivery := &Delivery{
		User:    (*userDelivery)(deliveryType),
		Article: (*articleDelivery)(deliveryType),
	}

	Route(e, delivery, mid)

	return delivery
}

func Route(e *echo.Echo, delivery *Delivery, mid *middlewares.CustomMiddleware) {
	e.POST("/register", delivery.User.RegisterUser)
	e.POST("/login", delivery.User.LoginUser)
	e.POST("/logout", delivery.User.LogoutUser, mid.JWT.ValidateJWT())

	// user
	user := e.Group("/user")
	{
		user.GET("", delivery.User.GetUser, mid.JWT.ValidateJWT())
		user.PUT("", delivery.User.UpdateUser, mid.JWT.ValidateJWT())
		user.DELETE("", delivery.User.DeleteUser, mid.JWT.ValidateJWT())
	}

	// article
	article := e.Group("/article")
	{
		article.GET("", delivery.Article.GetAllArticle, mid.JWT.ValidateJWT())
		article.POST("", delivery.Article.CreateArticle, mid.JWT.ValidateJWT())
	}
}
