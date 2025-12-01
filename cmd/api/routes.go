package main

import (
	"time"

	"github.com/gin-gonic/gin"

	"khalif-identify/internal/config"
	"khalif-identify/pkg/middleware"

)

func SetupRoutes(r *gin.Engine, app *App, cfg *config.Config) {
	r.Static("/uploads", "./uploads")

	globalLimitConfig := middleware.RateLimitConfig{
		Limit:  60,
		Window: time.Minute,
	}
	r.Use(middleware.RateLimit(app.RDB, globalLimitConfig))

	apiAdmin := r.Group("/api/admin")
	{
		apiAdmin.GET("/meta/countries", app.UserHandler.GetCountryCodes)
		apiAdmin.POST("/register", app.UserHandler.Register)

		loginLimit := middleware.RateLimitConfig{Limit: 5, Window: time.Minute}
		apiAdmin.POST("/login", middleware.RateLimit(app.RDB, loginLimit), app.UserHandler.Login)

		protectedAdmin := apiAdmin.Group("/")
		protectedAdmin.Use(middleware.AuthMiddleware(cfg.JWTSecret, app.RDB))
		{
			protectedAdmin.GET("/me", app.UserHandler.GetProfile)
			protectedAdmin.POST("/logout", app.UserHandler.Logout)
			protectedAdmin.POST("/profile/update", app.UserHandler.UpdateProfile)
			protectedAdmin.GET("/list", middleware.OnlyAdmin(), app.UserHandler.GetAll)
		}
	}

	apiUser := r.Group("/api/user")
	{
		apiUser.POST("/register", app.UserHandler.RegisterCustomer)
		
		loginLimit := middleware.RateLimitConfig{Limit: 5, Window: time.Minute}
		apiUser.POST("/login", middleware.RateLimit(app.RDB, loginLimit), app.UserHandler.Login)

		protectedUser := apiUser.Group("/")
		protectedUser.Use(middleware.AuthMiddleware(cfg.JWTSecret, app.RDB))
		{
			protectedUser.GET("/me", app.UserHandler.GetProfile)
			protectedUser.POST("/logout", app.UserHandler.Logout)
			protectedUser.POST("/profile/update", app.UserHandler.UpdateProfile)
		}
	}
}