package rest

import (
	"errors"
	"net/http"
	"strconv"
	"xis-data-aggregator/internal/repository"

	"xis-data-aggregator/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DataHandler struct {
	service *service.DataService
}

func NewDataHandler(service *service.DataService) *DataHandler {
	return &DataHandler{service: service}
}

// GetByID godoc
// @Summary      Get data by ID
// @Description  get data by UUID
// @Tags         data
// @Param        id   path      string  true  "Data ID"
// @Success      200  {object}  models.Data
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /data/{id} [get]
func (h *DataHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	data, err := h.service.GetByID(id)
	switch {
	case errors.Is(err, repository.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	case errors.Is(err, service.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// ListByTimeRange godoc
// @Summary      List data by time range
// @Description  get data by time range
// @Tags         data
// @Param        from  query     int64  true  "From timestamp"
// @Param        to    query     int64  true  "To timestamp"
// @Success      200  {array}   models.Data
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /data [get]
func (h *DataHandler) ListByTimeRange(c *gin.Context) {
	fromStr := c.Query("from")
	toStr := c.Query("to")
	from, err1 := strconv.ParseInt(fromStr, 10, 64)
	to, err2 := strconv.ParseInt(toStr, 10, 64)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from/to"})
		return
	}

	data, err := h.service.ListByPeriod(from, to)
	switch {
	case errors.Is(err, repository.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	case errors.Is(err, service.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, data)
}
