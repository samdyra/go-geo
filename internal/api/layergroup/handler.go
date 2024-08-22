package layergroup

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

func (h *Handler) CreateGroup(c *gin.Context) {
	var input LayerGroupCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrInvalidInput))
		return
	}

	username, _ := c.Get("username")
	err := h.service.CreateGroup(input, username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Group created successfully"})
}

func (h *Handler) AddLayerToGroup(c *gin.Context) {
	var input LayerToGroup
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrInvalidInput))
		return
	}

	username, _ := c.Get("username")
	err := h.service.AddLayerToGroup(input, username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Layer added to group successfully"})
}

func (h *Handler) GetGroupsWithLayers(c *gin.Context) {
	groups, err := h.service.GetGroupsWithLayers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
		return
	}

	c.JSON(http.StatusOK, groups)
}

func (h *Handler) RemoveLayerFromGroup(c *gin.Context) {
	layerID, err := strconv.ParseInt(c.Query("layer_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrInvalidInput))
		return
	}

	groupID, err := strconv.ParseInt(c.Query("group_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrInvalidInput))
		return
	}

	err = h.service.RemoveLayerFromGroup(layerID, groupID)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, errors.NewAPIError(err))
		default:
			c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Layer removed from group successfully"})
}

func (h *Handler) DeleteGroup(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrInvalidInput))
		return
	}

	err = h.service.DeleteGroup(groupID)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, errors.NewAPIError(err))
		default:
			c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group deleted successfully"})
}