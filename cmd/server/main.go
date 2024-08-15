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
	authHandler := api.NewHandler(authService)

	articleService := services.NewArticleService(db)
	articleHandler := api.NewArticleHandler(articleService)

	geoService := services.NewGeoService(db)
	geoHandler := api.NewGeoHandler(geoService)

	mvtService := services.NewMVTService(db)
	mvtHandler := api.NewMVTHandler(mvtService)

	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	// CORS middleware
	r.Use(cors.New(cors.Config{
		// @TODO: Change this to the actual frontend URL from env
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Public routes
	r.POST("/signup", authHandler.SignUp)
	r.POST("/signin", authHandler.SignIn)
	r.POST("/logout", authHandler.Logout)
	r.GET("/articles", articleHandler.GetArticles)
	r.GET("/articles/:id", articleHandler.GetArticle)
	r.GET("/mvt/:table_name/:z/:x/:y", mvtHandler.GetMVT)
	r.GET("/geo-data-list", geoHandler.GetGeoDataList)

	// Protected routes group
	protected := r.Group("/")
	protected.Use(middleware.JWTAuth())
	{
		articles := protected.Group("articles")
		{
			articles.POST("", articleHandler.CreateArticle)
			articles.PUT("/:id", articleHandler.UpdateArticle)
			articles.DELETE("/:id", articleHandler.DeleteArticle)
		}

		geo := protected.Group("geo")
		{
			geo.POST("/upload", geoHandler.UploadGeoData)
			geo.DELETE("/:table_name", geoHandler.DeleteGeoData)
		}
	}

	log.Printf("Starting server on :%s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}