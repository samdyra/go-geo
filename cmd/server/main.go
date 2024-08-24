package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/samdyra/go-geo/internal/api/article"
	"github.com/samdyra/go-geo/internal/api/geojson"
	"github.com/samdyra/go-geo/internal/api/layer"
	"github.com/samdyra/go-geo/internal/api/layergroup"
	"github.com/samdyra/go-geo/internal/api/mvt"
	"github.com/samdyra/go-geo/internal/api/report" // New import
	"github.com/samdyra/go-geo/internal/api/spatialdata"
	"github.com/samdyra/go-geo/internal/api/user"
	"github.com/samdyra/go-geo/internal/config"
	"github.com/samdyra/go-geo/internal/database"
	"github.com/samdyra/go-geo/internal/middleware"
)

func main() {
	cfg := config.Load()
	db := database.NewDB(cfg)
	
	authService := user.NewAuthService(db)
	authHandler := user.NewHandler(authService)

	articleService := article.NewArticleService(db)
	articleHandler := article.NewArticleHandler(articleService)

	spatialDataService := spatialdata.NewSpatialDataService(db)
	spatialDataHandler := spatialdata.NewSpatialDataHandler(spatialDataService)

	layerService := layer.NewService(db)
	layerHandler := layer.NewHandler(layerService)

	layerGroupService := layergroup.NewService(db)
	layerGroupHandler := layergroup.NewHandler(layerGroupService)

	mvtService := mvt.NewMVTService(db)
	mvtHandler := mvt.NewMVTHandler(mvtService)

	geoJSONService := geojson.NewGeoJSONService(db)
	geoJSONHandler := geojson.NewGeoJSONHandler(geoJSONService)

	reportService := report.NewReportService(db) 
	reportHandler := report.NewReportHandler(reportService)

	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	// CORS middleware
	r.Use(cors.New(cors.Config{
		// @TODO: Change this to the actual frontend URL from env
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174"},
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
	r.GET("/geojson/:table_name", geoJSONHandler.GetGeoJSON)
	r.GET("/layer-groups", layerGroupHandler.GetGroupsWithLayers)
	r.GET("/layers", layerHandler.GetFormattedLayers)
	r.POST("/reports", reportHandler.CreateReport)
	r.GET("/reports", reportHandler.GetReports)
	r.GET("/reports/:id", reportHandler.GetReport)

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

		spatialData := protected.Group("spatial-data")
		{
			spatialData.POST("", spatialDataHandler.CreateSpatialData)
			spatialData.DELETE("/:table_name", spatialDataHandler.DeleteSpatialData)
			spatialData.PUT("/:table_name", spatialDataHandler.EditSpatialData)
			spatialData.GET("", spatialDataHandler.GetSpatialDataList)
		}

		layers := protected.Group("layers")
		{
			layers.POST("", layerHandler.CreateLayer)
			layers.PUT("/:id", layerHandler.UpdateLayer)
			layers.DELETE("/:id", layerHandler.DeleteLayer)
		}

		layerGroups := protected.Group("layer-groups")
		{
			layerGroups.POST("", layerGroupHandler.CreateGroup)
			layerGroups.POST("/add-layer", layerGroupHandler.AddLayerToGroup)
			layerGroups.DELETE("/remove-layer", layerGroupHandler.RemoveLayerFromGroup)
			layerGroups.DELETE("/:id", layerGroupHandler.DeleteGroup)
		}

		reports := protected.Group("reports")  // New protected group for reports
		{
			reports.PUT("/:id", reportHandler.UpdateReport)
			reports.DELETE("/:id", reportHandler.DeleteReport)
		}
	}

	log.Printf("Starting server on :%s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}