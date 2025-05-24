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

// RegisterRoutes registers the API routes in the HTTP server.
func RegisterRoutes(e *echo.Echo) {
	e.GET("/", handleIndex)
	e.POST("/calculate", handleCalculation)
}

func handleIndex(c echo.Context) error {
	dir, err := os.Getwd()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get current directory")
	}

	return c.File(filepath.Join(dir, "public/index.html"))
}

func handleCalculation(c echo.Context) error {
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

	result, err := pack.CalculatePacking(req.OrderedItems, req.BoxSizes...)
	if err != nil {
		zap.L().Error("failed to calculate packing", zap.Error(err))

		return c.String(http.StatusInternalServerError, "Failed to calculate packing")
	}

	return c.JSON(http.StatusOK, result)
}

type calculateRequest struct {
	OrderedItems int   `json:"orderedItems" validate:"required,gt=0,lte=9007199254740991"`
	BoxSizes     []int `json:"boxSizes" validate:"required,unique,dive,gt=0,lte=9007199254740991"`
}
