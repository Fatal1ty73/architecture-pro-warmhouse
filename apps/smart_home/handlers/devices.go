package handlers

import (
	"io"
	"net/http"

	"smarthome/services"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	DeviceService *services.DeviceHandleService
}

func NewDeviceHandler(deviceService *services.DeviceHandleService) *DeviceHandler {
	return &DeviceHandler{DeviceService: deviceService}
}

func (h *DeviceHandler) RegisterRoutes(router *gin.RouterGroup) {
	devices := router.Group("/devices")
	{
		devices.GET("", h.ListDevices)
		devices.GET("/:id", h.GetDevice)
		devices.POST("", h.CreateDevice)
		devices.PATCH("/:id", h.UpdateDevice)
		devices.DELETE("/:id", h.DeleteDevice)
		devices.POST("/:id/execute", h.ExecuteDevice)
	}
}

func (h *DeviceHandler) ListDevices(c *gin.Context) {
	h.forward(c, http.MethodGet, "/api/v1/devices")
}

func (h *DeviceHandler) GetDevice(c *gin.Context) {
	h.forward(c, http.MethodGet, "/api/v1/devices/"+c.Param("id"))
}

func (h *DeviceHandler) CreateDevice(c *gin.Context) {
	h.forward(c, http.MethodPost, "/api/v1/devices")
}

func (h *DeviceHandler) UpdateDevice(c *gin.Context) {
	h.forward(c, http.MethodPatch, "/api/v1/devices/"+c.Param("id"))
}

func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	h.forward(c, http.MethodDelete, "/api/v1/devices/"+c.Param("id"))
}

func (h *DeviceHandler) ExecuteDevice(c *gin.Context) {
	h.forward(c, http.MethodPost, "/api/v1/devices/"+c.Param("id")+"/execute")
}

func (h *DeviceHandler) forward(c *gin.Context, method, path string) {
	if h.DeviceService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "device handle service is not configured"})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	status, payload, err := h.DeviceService.Forward(method, path, c.Request.URL.Query(), body)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	if len(payload) == 0 {
		c.Status(status)
		return
	}
	c.Data(status, "application/json", payload)
}
