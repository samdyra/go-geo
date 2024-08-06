package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	articleService := services.NewArticleService(db)
	authHandler := api.NewHandler(authService)
	articleHandler := api.NewArticleHandler(articleService)
	geoService := services.NewGeoService(db)
	geoHandler := api.NewGeoHandler(geoService)

	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Auth routes
	r.POST("/signup", authHandler.SignUp)
	r.POST("/signin", authHandler.SignIn)
	r.POST("/logout", authHandler.Logout)

	// Article routes
	r.GET("/articles", articleHandler.GetArticles)
	r.GET("/articles/:id", articleHandler.GetArticle)

	// Geo routes
	r.POST("/geo/upload", geoHandler.UploadGeoData)
	r.DELETE("/geo/:table_name", geoHandler.DeleteGeoData)

	// Protected article routes
	protected := r.Group("/articles")
	protected.Use(middleware.JWTAuth())
	{
		protected.POST("", articleHandler.CreateArticle)
		protected.PUT("/:id", articleHandler.UpdateArticle)
		protected.DELETE("/:id", articleHandler.DeleteArticle)
	}

	log.Printf("Starting server on :%s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}