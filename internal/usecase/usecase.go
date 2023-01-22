package usecase

import (
	"github.com/haikalvidya/go-article/config"
	"github.com/haikalvidya/go-article/internal/middlewares"
	"github.com/haikalvidya/go-article/internal/repository"

	"github.com/go-redis/redis"
)

type Usecase struct {
	User    IUserUsecase
	Article IArticleUsecase
}

type usecaseType struct {
	Repo        *repository.Repository
	Middleware  *middlewares.CustomMiddleware
	RedisClient *redis.Client
	ServerInfo  *config.ServerConfig
}

func NewUsecase(repo *repository.Repository, mid *middlewares.CustomMiddleware, redis *redis.Client, serverInfo *config.ServerConfig) *Usecase {
	usc := &usecaseType{Repo: repo, Middleware: mid, RedisClient: redis, ServerInfo: serverInfo}

	return &Usecase{
		User:    (*userUsecase)(usc),
		Article: (*articleUsecase)(usc),
	}
}
