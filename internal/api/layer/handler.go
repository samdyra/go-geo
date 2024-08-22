package layer

import (
	"net/http"
	"strconv"
	"strings"

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

    // Validate coordinate
    if len(input.Coordinate) != 2 {
        c.JSON(http.StatusBadRequest, errors.ErrInternalServer)
        return
    }

    username, exists := c.Get("username")
    if !exists {
        c.JSON(http.StatusUnauthorized, errors.NewAPIError(errors.ErrUnauthorized))
        return
    }

    err := h.service.CreateLayer(input, username.(string))
    if err != nil {
        c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Layer created successfully"})
}

func (h *Handler) GetFormattedLayers(c *gin.Context) {


	idParam := c.Query("id")

	if idParam == "" {

		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide 'id' parameter. Use 'id=*' to fetch all layers"})
		return
	}

	var ids []int64
	var err error

	if idParam == "*" {

		layers, err := h.service.GetAllFormattedLayers()
		if err != nil {

			c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
			return
		}

		c.JSON(http.StatusOK, layers)
		return
	}

	// Parse provided IDs
	idStrings := strings.Split(idParam, ",")
	for _, idStr := range idStrings {
		id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
		if err != nil {

			c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrInvalidInput))
			return
		}
		ids = append(ids, id)
	}


	layers, err := h.service.GetFormattedLayers(ids)
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


    username, exists := c.Get("username")
    if !exists {
        c.JSON(http.StatusUnauthorized, errors.NewAPIError(errors.ErrUnauthorized))
        return
    }


    err = h.service.UpdateLayer(id, input, username.(string))
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