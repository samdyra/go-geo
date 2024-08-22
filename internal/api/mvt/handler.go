package mvt

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/samdyra/go-geo/internal/utils/errors"
)

type MVTHandler struct {
	mvtService *MVTService
}

func NewMVTHandler(mvtService *MVTService) *MVTHandler {
	return &MVTHandler{mvtService: mvtService}
}

func (h *MVTHandler) GetMVT(c *gin.Context) {
	tableName := c.Param("table_name")
	z, _ := strconv.Atoi(c.Param("z"))
	x, _ := strconv.Atoi(c.Param("x"))
	y, _ := strconv.Atoi(c.Param("y"))

	mvt, err := h.mvtService.GenerateMVT(tableName, z, x, y)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
		return
	}

	c.Header("Content-Type", "application/x-protobuf")
	c.Data(http.StatusOK, "application/x-protobuf", mvt)
}