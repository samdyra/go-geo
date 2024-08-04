package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"github.com/samdyra/go-geo/internal/api"
	"github.com/samdyra/go-geo/internal/config"
	"github.com/samdyra/go-geo/internal/database"
	"github.com/samdyra/go-geo/internal/middleware"
	"github.com/samdyra/go-geo/internal/services"
)

func main() {
	cfg := config.Load()
	db := database.NewDB(cfg)
	authService := services.NewAuthService(db)
	handler := api.NewHandler(authService)

	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rate limiting middleware
	limiter := middleware.NewIPRateLimiter(rate.Limit(1), 5) // 1 request per second with burst of 5
	r.Use(middleware.RateLimitMiddleware(limiter))

	// Public routes
	r.POST("/signup", handler.SignUp)
	r.POST("/signin", handler.SignIn)
	r.POST("/logout", handler.Logout)

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth())
	{
		protected.GET("/protected", handler.ProtectedRoute)
	}

	log.Printf("Starting server on :%s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}