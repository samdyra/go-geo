// report/handler.go

package report

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/samdyra/go-geo/internal/utils/errors"
)

type ReportHandler struct {
	reportService *ReportService
}

func NewReportHandler(reportService *ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

func (h *ReportHandler) GetReports(c *gin.Context) {
	reports, err := h.reportService.GetReports()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
		return
	}
	c.JSON(http.StatusOK, reports)
}

func (h *ReportHandler) GetReport(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrInvalidInput))
		return
	}

	report, err := h.reportService.GetReportByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, errors.NewAPIError(err))
		return
	}
	c.JSON(http.StatusOK, report)
}

func (h *ReportHandler) CreateReport(c *gin.Context) {
    var input CreateReportInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrInvalidInput))
        return
    }

    if err := input.Validate(); err != nil {
        c.JSON(http.StatusBadRequest, errors.NewAPIError(err))
        return
    }

    report, err := h.reportService.CreateReport(input)
    if err != nil {
        c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
        return
    }

    c.JSON(http.StatusCreated, report)
}

func (h *ReportHandler) UpdateReport(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrInvalidInput))
		return
	}

	var input UpdateReportInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrInvalidInput))
		return
	}

	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(err))
		return
	}

	report, err := h.reportService.UpdateReport(id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
		return
	}

	c.JSON(http.StatusOK, report)
}

func (h *ReportHandler) DeleteReport(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrInvalidInput))
		return
	}

	err = h.reportService.DeleteReport(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Report deleted successfully"})
}