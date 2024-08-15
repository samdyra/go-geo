package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samdyra/go-geo/internal/models"
	"github.com/samdyra/go-geo/internal/services"
	"github.com/samdyra/go-geo/internal/utils/errors"
)

type GeoHandler struct {
    geoService *services.GeoService
}

func NewGeoHandler(geoService *services.GeoService) *GeoHandler {
    return &GeoHandler{geoService: geoService}
}

func (h *GeoHandler) GetGeoDataList(c *gin.Context) {
    formattedData, err := h.geoService.GetFormattedGeoData()
    if err != nil {
        c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
        return
    }

    c.JSON(http.StatusOK, formattedData)
}

func (h *GeoHandler) UploadGeoData(c *gin.Context) {
    var input models.GeoDataUpload
    if err := c.ShouldBind(&input); err != nil {
        c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrInvalidInput))
        return
    }

    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrInvalidInput))
        return
    }
    
    if len(input.Coordinate) != 2 {
        c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrInvalidInput))
        return
    }

    openedFile, err := file.Open()
    if err != nil {
        c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
        return
    }
    defer openedFile.Close()

    username, exists := c.Get("username")
    if !exists {
        c.JSON(http.StatusUnauthorized, errors.NewAPIError(errors.ErrUnauthorized))
        return
    }

    err = h.geoService.CreateGeoData(input, openedFile, username.(string))
    if err != nil {
        c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Geo data uploaded successfully"})
}


func (h *GeoHandler) DeleteGeoData(c *gin.Context) {
    tableName := c.Param("table_name")

    err := h.geoService.DeleteGeoData(tableName)
    if err != nil {
        c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Geo data deleted successfully"})
}