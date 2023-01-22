package middlewares

import (
	"github.com/haikalvidya/go-article/config"
	"github.com/haikalvidya/go-article/pkg/logger"
	"github.com/haikalvidya/go-article/pkg/middleware"

	"github.com/labstack/echo/v4"
)

type jwtImpl interface {
	GenerateToken(userId []byte) (string, error)
	ValidateJWT() echo.MiddlewareFunc
	GetJWTClaims(c echo.Context) map[string]interface{}
	GetUserIdFromJwt(c echo.Context) string
}

type CustomMiddleware struct {
	JWT            jwtImpl
	Logger         logger.Logger
	InternalAccess internalconnection
}

type internalconnection interface {
	ValidateInternalAccess(next echo.HandlerFunc) echo.HandlerFunc
}

type customMiddleware struct {
	Config *config.Config
}

func New(cfg *config.Config) *CustomMiddleware {

	jwt := middleware.NewJwt(cfg.JWT.AccessTokenExpiredHour, cfg.JWT.Secret)

	logger := logger.NewApiLogger(cfg)

	internalconnection := middleware.NewInternalAccess(cfg.Server.InternalAccessKey)

	return &CustomMiddleware{
		JWT:            jwt,
		Logger:         logger,
		InternalAccess: internalconnection,
	}
}
