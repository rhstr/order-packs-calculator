package api

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/rhstr/order-packs-calculator/internal/pack"
)

// HTTPHandler represents a handler for HTTP requests.
type HTTPHandler interface {
	RegisterRoutes(e *echo.Echo)
}

type handler struct {
	cache pack.Cache
}

// NewHandler creates a new instance of HTTPHandler with the provided cache.
func NewHandler(cache pack.Cache) HTTPHandler {
	return &handler{
		cache: cache,
	}
}

// RegisterRoutes registers the API routes in the HTTP server.
func (h *handler) RegisterRoutes(e *echo.Echo) {
	e.GET("/", h.handleIndex)
	e.POST("/calculate", h.handleCalculation)
}

func (h *handler) handleIndex(c echo.Context) error {
	dir, err := os.Getwd()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get current directory")
	}

	return c.File(filepath.Join(dir, "public/index.html"))
}

func (h *handler) handleCalculation(c echo.Context) error {
	var req calculateRequest
	if err := c.Bind(&req); err != nil {
		zap.L().Error("failed to bind request", zap.Error(err))

		return c.String(http.StatusBadRequest, "Invalid request")
	}

	err := validator.New().Struct(&req)
	if err != nil {
		zap.L().Error("validation failed", zap.Error(err))

		return c.String(http.StatusBadRequest, "Invalid request")
	}

	result := h.cache.Get(req.OrderedItems, req.BoxSizes...)
	if result != nil {
		return c.JSON(http.StatusOK, result)
	}

	result, err = pack.CalculatePacking(req.OrderedItems, req.BoxSizes...)
	if err != nil {
		zap.L().Error("failed to calculate packing", zap.Error(err))

		return c.String(http.StatusInternalServerError, "Failed to calculate packing")
	}

	h.cache.Set(req.OrderedItems, req.BoxSizes, result)

	return c.JSON(http.StatusOK, result)
}

// calculateRequest represents the request body for the packing calculation.
// Ordered items and box sizes cannot be greater than (2^53 - 1) due to JS limitations.
type calculateRequest struct {
	OrderedItems int   `json:"orderedItems" validate:"required,gt=0,lte=9007199254740991"`
	BoxSizes     []int `json:"boxSizes" validate:"required,unique,dive,gt=0,lte=9007199254740991"`
}
