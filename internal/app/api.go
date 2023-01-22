package app

import (
	"context"
	"log"

	"github.com/haikalvidya/go-article/internal/delivery"
	"github.com/haikalvidya/go-article/internal/middlewares"
	"github.com/haikalvidya/go-article/internal/repository"
	"github.com/haikalvidya/go-article/internal/usecase"

	"github.com/haikalvidya/go-article/pkg"

	"github.com/haikalvidya/go-article/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type httpApp struct {
	base
	router     *echo.Echo
	usecase    *usecase.Usecase
	repo       *repository.Repository
	delivery   *delivery.Delivery
	middleware *middlewares.CustomMiddleware
	signalHttp *pkg.GracefullShutdown
}

func (a *httpApp) Init() (err error) {
	err = a.initConfig()
	if err != nil {
		return
	}
	a.repo = repository.NewRepository(a.db)
	a.middleware = middlewares.New(a.config)
	a.usecase = usecase.NewUsecase(a.repo, a.middleware, a.redis, &a.config.Server)

	e := echo.New()

	e.Validator = &utils.CustomValidator{Validator: validator.New()}

	e.Use(utils.RateLimit())
	e.Use(middleware.SecureWithConfig(middleware.DefaultSecureConfig))
	e.IPExtractor = echo.ExtractIPDirect()
	a.router = e

	a.delivery = delivery.NewDelivery(a.router, a.usecase, a.middleware)
	a.signalHttp = pkg.NewGracefullShutdown()
	return
}

func (a *httpApp) Run() (err error) {
	go func() {
		a.signalHttp.Wait()

		log.Println("Shutting down the service!")
		if err := a.router.Shutdown(context.Background()); err != nil {
			log.Printf("Error in shutdown the service: %v.", err)
		}
	}()

	log.Println("Press Ctrl + C to exit the service!")

	err = a.router.Start(a.config.Server.Address)
	return
}

func (a *httpApp) Close() (err error) {

	// close base config
	a.closeConfig()

	return
}
