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

// UploadGeoData godoc
// @Summary Upload geospatial data
// @Description Upload a GeoJSON file and create a new geo data entry
// @Tags geo
// @Accept multipart/form-data
// @Produce json
// @Param table_name formData string true "Name for the new table"
// @Param type formData string true "Type of geometry"
// @Param coordinate formData string false "Centroid coordinate"
// @Param file formData file true "GeoJSON file"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} errors.APIError
// @Failure 500 {object} errors.APIError
// @Router /geo/upload [post]
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

    openedFile, err := file.Open()
    if err != nil {
        c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
        return
    }
    defer openedFile.Close()

    err = h.geoService.CreateGeoData(input, openedFile)
    if err != nil {
        c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Geo data uploaded successfully"})
}

// DeleteGeoData godoc
// @Summary Delete geospatial data
// @Description Delete a geo data entry and its corresponding table
// @Tags geo
// @Produce json
// @Param table_name path string true "Table name to delete"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} errors.APIError
// @Failure 500 {object} errors.APIError
// @Router /geo/{table_name} [delete]
func (h *GeoHandler) DeleteGeoData(c *gin.Context) {
    tableName := c.Param("table_name")

    err := h.geoService.DeleteGeoData(tableName)
    if err != nil {
        c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Geo data deleted successfully"})
}