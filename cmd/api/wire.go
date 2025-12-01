//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"khalif-identify/internal/config"
	"khalif-identify/internal/domain"
	"khalif-identify/internal/handler"
	"khalif-identify/internal/repository"
	"khalif-identify/internal/usecase"
	"khalif-identify/pkg/utils"

)

// 1. Buat Struct Penampung (Container)
type App struct {
	DB          *gorm.DB
	RDB         *redis.Client
	UserHandler *handler.UserHandler
}

// 2. Provider untuk membuat Struct App
func NewApp(db *gorm.DB, rdb *redis.Client, h *handler.UserHandler) *App {
	return &App{
		DB:          db,
		RDB:         rdb,
		UserHandler: h,
	}
}

// 3. Update InitializeApp agar me-return *App saja
func InitializeApp() (*App, error) {
	wire.Build(
		config.LoadConfig,
		ProvideDB,
		ProvideRedis,
		ProvideAzureUploader,
		ProvideJWTSecret,

		repository.NewUserRepository,
		wire.Bind(new(domain.UserRepository), new(*repository.UserRepo)),

		repository.NewCacheRepository,
		wire.Bind(new(domain.CacheRepository), new(*repository.RedisRepo)),

		NewUserUseCaseWire,
		handler.NewUserHandler,

		// Masukkan Provider App Baru
		NewApp,
	)
	return &App{}, nil
}

func NewUserUseCaseWire(
	repo domain.UserRepository,
	cache domain.CacheRepository,
	uploader *utils.AzureUploader,
	secret JWTSecret,
) usecase.UserUseCase {
	return usecase.NewUserUseCase(repo, cache, uploader, string(secret))
}