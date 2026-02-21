package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"sensor-manager-service/internal/db"
	"sensor-manager-service/internal/model"
	"sensor-manager-service/internal/service"

	"github.com/gin-gonic/gin"
)

type SensorHandler struct {
	repo       *db.Repository
	tempClient *service.TemperatureClient
}

func NewSensorHandler(repo *db.Repository, tempClient *service.TemperatureClient) *SensorHandler {
	return &SensorHandler{repo: repo, tempClient: tempClient}
}

func (h *SensorHandler) Register(router *gin.RouterGroup) {
	sensors := router.Group("/sensors")
	sensors.GET("", h.ListSensors)
	sensors.GET("/:id", h.GetSensor)
	sensors.POST("", h.CreateSensor)
	sensors.PATCH("/:id", h.UpdateSensor)
	sensors.DELETE("/:id", h.DeleteSensor)
}

func (h *SensorHandler) ListSensors(c *gin.Context) {
	limit, offset := parsePagination(c)
	houseID := parseInt64Query(c, "house_id")
	typeID := parseInt64Query(c, "type_id")
	status := parseStringQuery(c, "status")

	items, total, err := h.repo.ListSensors(c.Request.Context(), limit, offset, houseID, typeID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errBody("INTERNAL_ERROR", err.Error()))
		return
	}

	for i := range items {
		h.refreshTelemetry(c.Request.Context(), &items[i])
	}

	c.JSON(http.StatusOK, model.SensorList{Items: items, Limit: limit, Offset: offset, Total: total})
}

func (h *SensorHandler) GetSensor(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errBody("BAD_REQUEST", "invalid sensor id"))
		return
	}

	s, err := h.repo.GetSensor(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, errBody("NOT_FOUND", "sensor not found"))
		return
	}

	h.refreshTelemetry(c.Request.Context(), s)
	c.JSON(http.StatusOK, s)
}

func (h *SensorHandler) CreateSensor(c *gin.Context) {
	var in model.SensorCreate
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, errBody("BAD_REQUEST", err.Error()))
		return
	}

	s, err := h.repo.CreateSensor(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, errBody("BAD_REQUEST", err.Error()))
		return
	}
	c.JSON(http.StatusCreated, s)
}

func (h *SensorHandler) UpdateSensor(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errBody("BAD_REQUEST", "invalid sensor id"))
		return
	}

	var in model.SensorUpdate
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, errBody("BAD_REQUEST", err.Error()))
		return
	}

	s, err := h.repo.UpdateSensor(c.Request.Context(), id, in)
	if err != nil {
		c.JSON(http.StatusBadRequest, errBody("BAD_REQUEST", err.Error()))
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *SensorHandler) DeleteSensor(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errBody("BAD_REQUEST", "invalid sensor id"))
		return
	}

	if err := h.repo.DeleteSensor(c.Request.Context(), id); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, errBody("NOT_FOUND", "sensor not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, errBody("INTERNAL_ERROR", err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *SensorHandler) refreshTelemetry(ctx context.Context, s *model.Sensor) {
	temp, err := h.tempClient.FetchBySensorID(s.ID)
	if err != nil {
		temp, err = h.tempClient.FetchByLocation(s.Location)
	}
	if err != nil {
		return
	}
	s.Value = temp.Value
	s.Unit = temp.Unit
	s.Status = temp.Status
	s.LastUpdated = temp.Timestamp
	_ = h.repo.ApplyTelemetry(ctx, s.ID, temp.Value, temp.Unit, temp.Status, temp.Timestamp)
}

func parsePagination(c *gin.Context) (int, int) {
	limit := 50
	offset := 0
	if raw := c.Query("limit"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v > 0 && v <= 200 {
			limit = v
		}
	}
	if raw := c.Query("offset"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v >= 0 {
			offset = v
		}
	}
	return limit, offset
}

func parseInt64Query(c *gin.Context, key string) *int64 {
	v := c.Query(key)
	if v == "" {
		return nil
	}
	parsed, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil
	}
	return &parsed
}

func parseStringQuery(c *gin.Context, key string) *string {
	v := c.Query(key)
	if v == "" {
		return nil
	}
	return &v
}

func errBody(code, msg string) gin.H {
	return gin.H{"code": code, "message": msg}
}
