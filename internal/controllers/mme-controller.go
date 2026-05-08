package controllers

import (
	"4g-proxy/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// MMEController handles HTTP requests for proxy management
type MMEController struct {
	state *models.ProxyState
}

// NewMMEController creates a new MMEController
func NewMMEController(state *models.ProxyState) *MMEController {
	return &MMEController{
		state: state,
	}
}

// GetStatus returns the current proxy status
func (c *MMEController) GetStatus(ctx *gin.Context) {
	stats := c.state.GetStats()
	dropFlags := c.state.DropFlags.GetAll()
	delayConfig := c.state.DelayConfig.GetAll()

	ctx.JSON(http.StatusOK, gin.H{
		"status":      "running",
		"stats":       stats,
		"dropFlags":   dropFlags,
		"delayConfig": delayConfig,
	})
}

// GetStats returns proxy statistics
func (c *MMEController) GetStats(ctx *gin.Context) {
	stats := c.state.GetStats()
	ctx.JSON(http.StatusOK, stats)
}

// GetDropFlags returns current drop flags
func (c *MMEController) GetDropFlags(ctx *gin.Context) {
	flags := c.state.DropFlags.GetAll()
	ctx.JSON(http.StatusOK, flags)
}

// SetDropFlag sets a drop flag
func (c *MMEController) SetDropFlag(ctx *gin.Context) {
	signalType := ctx.Param("signalType")

	var body struct {
		Drop bool `json:"drop"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if !c.state.DropFlags.SetDropByName(signalType, body.Drop) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Unknown signal type: " + signalType,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Drop flag updated",
		"signalType": signalType,
		"drop":       body.Drop,
	})
}

// SetDropFlags sets multiple drop flags at once
func (c *MMEController) SetDropFlags(ctx *gin.Context) {
	var body map[string]bool

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	updated := make(map[string]bool)
	for signalType, drop := range body {
		if c.state.DropFlags.SetDropByName(signalType, drop) {
			updated[signalType] = drop
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Drop flags updated",
		"updated": updated,
	})
}

// ResetDropFlags resets all drop flags
func (c *MMEController) ResetDropFlags(ctx *gin.Context) {
	c.state.DropFlags.Reset()
	ctx.JSON(http.StatusOK, gin.H{
		"message": "All drop flags reset",
	})
}

// DropAttach enables dropping of Attach messages
func (c *MMEController) DropAttach(ctx *gin.Context) {
	c.state.DropFlags.SetDropByName("attach", true)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Attach messages will be dropped",
	})
}

// AllowAttach disables dropping of Attach messages
func (c *MMEController) AllowAttach(ctx *gin.Context) {
	c.state.DropFlags.SetDropByName("attach", false)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Attach messages will be forwarded",
	})
}

// DropDetach enables dropping of Detach messages
func (c *MMEController) DropDetach(ctx *gin.Context) {
	c.state.DropFlags.SetDropByName("detach", true)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Detach messages will be dropped",
	})
}

// AllowDetach disables dropping of Detach messages
func (c *MMEController) AllowDetach(ctx *gin.Context) {
	c.state.DropFlags.SetDropByName("detach", false)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Detach messages will be forwarded",
	})
}

// DropTAU enables dropping of TAU messages
func (c *MMEController) DropTAU(ctx *gin.Context) {
	c.state.DropFlags.SetDropByName("tau", true)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "TAU messages will be dropped",
	})
}

// AllowTAU disables dropping of TAU messages
func (c *MMEController) AllowTAU(ctx *gin.Context) {
	c.state.DropFlags.SetDropByName("tau", false)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "TAU messages will be forwarded",
	})
}

// DropServiceRequest enables dropping of Service Request messages
func (c *MMEController) DropServiceRequest(ctx *gin.Context) {
	c.state.DropFlags.SetDropByName("serviceRequest", true)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Service Request messages will be dropped",
	})
}

// AllowServiceRequest disables dropping of Service Request messages
func (c *MMEController) AllowServiceRequest(ctx *gin.Context) {
	c.state.DropFlags.SetDropByName("serviceRequest", false)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Service Request messages will be forwarded",
	})
}

// DropUEContextRelease enables dropping of UE Context Release messages
func (c *MMEController) DropUEContextRelease(ctx *gin.Context) {
	c.state.DropFlags.SetDropByName("ueContextRelease", true)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "UE Context Release messages will be dropped",
	})
}

// AllowUEContextRelease disables dropping of UE Context Release messages
func (c *MMEController) AllowUEContextRelease(ctx *gin.Context) {
	c.state.DropFlags.SetDropByName("ueContextRelease", false)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "UE Context Release messages will be forwarded",
	})
}

// Health returns health check status
func (c *MMEController) Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}

// GetDelayConfig returns current delay configuration
func (c *MMEController) GetDelayConfig(ctx *gin.Context) {
	delays := c.state.DelayConfig.GetAll()
	ctx.JSON(http.StatusOK, delays)
}

// SetDelay sets a delay for a signal type
func (c *MMEController) SetDelay(ctx *gin.Context) {
	signalType := ctx.Param("signalType")

	var body struct {
		DelayMs int64 `json:"delayMs"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if !c.state.DelayConfig.SetDelayByName(signalType, body.DelayMs) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Unknown signal type: " + signalType,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Delay updated",
		"signalType": signalType,
		"delayMs":    body.DelayMs,
	})
}

// SetDelays sets multiple delays at once
func (c *MMEController) SetDelays(ctx *gin.Context) {
	var body map[string]int64

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	updated := make(map[string]int64)
	for signalType, delayMs := range body {
		if c.state.DelayConfig.SetDelayByName(signalType, delayMs) {
			updated[signalType] = delayMs
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Delays updated",
		"updated": updated,
	})
}

// ResetDelays resets all delays to zero
func (c *MMEController) ResetDelays(ctx *gin.Context) {
	c.state.DelayConfig.Reset()
	ctx.JSON(http.StatusOK, gin.H{
		"message": "All delays reset to zero",
	})
}
