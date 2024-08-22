package spatialdata

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/samdyra/go-geo/internal/utils/errors"
)

type SpatialDataHandler struct {
    spatialDataService *SpatialDataService
}

func NewSpatialDataHandler(spatialDataService *SpatialDataService) *SpatialDataHandler {
    return &SpatialDataHandler{spatialDataService: spatialDataService}
}

func (h *SpatialDataHandler) CreateSpatialData(c *gin.Context) {
    var input SpatialDataCreate
    if err := c.ShouldBindWith(&input, binding.Form); err != nil {
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

    username, exists := c.Get("username")
    if !exists {
        c.JSON(http.StatusUnauthorized, errors.NewAPIError(errors.ErrUnauthorized))
        return
    }

    err = h.spatialDataService.CreateSpatialData(input, openedFile, username.(string))
    if err != nil {
        c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Spatial data created successfully"})
}

func (h *SpatialDataHandler) GetSpatialDataList(c *gin.Context) {
    spatialDataList, err := h.spatialDataService.GetSpatialDataList()
    if err != nil {
        c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
        return
    }

    c.JSON(http.StatusOK, spatialDataList)
}

func (h *SpatialDataHandler) DeleteSpatialData(c *gin.Context) {
    tableName := c.Param("table_name")

    err := h.spatialDataService.DeleteSpatialData(tableName)
    if err != nil {
        switch err {
        case errors.ErrNotFound:
            c.JSON(http.StatusNotFound, errors.NewAPIError(err))
        default:
            c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Spatial data deleted successfully"})
}

func (h *SpatialDataHandler) EditSpatialData(c *gin.Context) {
    var input SpatialDataEdit
    if err := c.ShouldBind(&input); err != nil {
        c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrInvalidInput))
        return
    }

    oldTableName := c.Param("table_name")
    username, _ := c.Get("username")

    var file io.Reader
    uploadedFile, err := c.FormFile("file")
    if err == nil {
        openedFile, err := uploadedFile.Open()
        if err != nil {
            c.JSON(http.StatusInternalServerError, errors.NewAPIError(errors.ErrInternalServer))
            return
        }
        defer openedFile.Close()
        file = openedFile
    }

    err = h.spatialDataService.EditSpatialData(oldTableName, input, file, username.(string))
    if err != nil {
        switch err {
        case errors.ErrNotFound:
            c.JSON(http.StatusNotFound, errors.NewAPIError(err))
        case errors.ErrInvalidInput:
            c.JSON(http.StatusBadRequest, errors.NewAPIError(err))
        default:
            c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Spatial data updated successfully"})
}
