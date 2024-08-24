package geojson

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samdyra/go-geo/internal/utils/errors"
)

type GeoJSONHandler struct {
	geojsonService *GeoJSONService
}

func NewGeoJSONHandler(geojsonService *GeoJSONService) *GeoJSONHandler {
	return &GeoJSONHandler{geojsonService: geojsonService}
}

func (h *GeoJSONHandler) GetGeoJSON(c *gin.Context) {
	tableName := c.Param("table_name")

	geojson, err := h.geojsonService.GenerateGeoJSON(tableName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
		return
	}

	c.Header("Content-Type", "application/json")
	c.Data(http.StatusOK, "application/json", geojson)
}

