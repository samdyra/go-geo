package layer

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/samdyra/go-geo/internal/utils/errors"
)

type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) CreateLayer(c *gin.Context) {
    var input LayerCreate
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrInvalidInput))
        return
    }

    userID, _ := c.Get("user_id")
    err := h.service.CreateLayer(input, userID.(int64))
    if err != nil {
        c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Layer created successfully"})
}

func (h *Handler) GetFormattedLayers(c *gin.Context) {
    layers, err := h.service.GetFormattedLayers()
    if err != nil {
        switch err {
        case errors.ErrNotFound:
            c.JSON(http.StatusNotFound, errors.NewAPIError(err))
        default:
            c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
        }
        return
    }

    c.JSON(http.StatusOK, layers)
}

func (h *Handler) UpdateLayer(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrInvalidInput))
        return
    }

    var input LayerUpdate
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrInvalidInput))
        return
    }

    userID, _ := c.Get("user_id")
    err = h.service.UpdateLayer(id, input, userID.(int64))
    if err != nil {
        switch err {
        case errors.ErrNotFound:
            c.JSON(http.StatusNotFound, errors.NewAPIError(err))
        default:
            c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Layer updated successfully"})
}

func (h *Handler) DeleteLayer(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrInvalidInput))
        return
    }

    err = h.service.DeleteLayer(id)
    if err != nil {
        switch err {
        case errors.ErrNotFound:
            c.JSON(http.StatusNotFound, errors.NewAPIError(err))
        default:
            c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Layer deleted successfully"})
}