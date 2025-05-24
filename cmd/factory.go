package cmd

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"github.com/rhstr/order-packs-calculator/internal/api"
	"github.com/rhstr/order-packs-calculator/internal/config"
	"github.com/rhstr/order-packs-calculator/internal/pack"
)

var serviceContainer = new(diContainer)

// diContainer is a dependency injection container for the application.
type diContainer struct {
	cfg config.Config

	httpServer  *echo.Echo
	httpHandler api.HTTPHandler

	packCalculator pack.Calculator
	packCache      pack.Cache
}

func (c *diContainer) getConfig() config.Config {
	if c.cfg == (config.Config{}) {
		cfg, err := config.New()
		if err != nil {
			zap.L().Fatal("failed to create config", zap.Error(err))
		}

		c.cfg = cfg
	}

	return c.cfg
}

func (c *diContainer) getHTTPServer() *echo.Echo {
	if c.httpServer == nil {
		c.httpServer = echo.New()

		corsCfg := middleware.DefaultCORSConfig
		corsCfg.AllowOrigins = []string{c.getConfig().CORSOrigin}

		c.httpServer.Use(middleware.CORSWithConfig(corsCfg))
		c.httpServer.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogLatency:       true,
			LogProtocol:      true,
			LogRemoteIP:      true,
			LogMethod:        true,
			LogURI:           true,
			LogRequestID:     true,
			LogReferer:       true,
			LogUserAgent:     true,
			LogStatus:        true,
			LogError:         true,
			LogContentLength: true,
			LogResponseSize:  true,
			HandleError:      true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				zap.L().Info("request",
					zap.String("protocol", v.Protocol),
					zap.String("remote_ip", v.RemoteIP),
					zap.String("method", v.Method),
					zap.String("uri", v.URI),
					zap.String("request_id", v.RequestID),
					zap.String("referer", v.Referer),
					zap.String("user_agent", v.UserAgent),
					zap.Int("status", v.Status),
					zap.Error(v.Error),
					zap.String("content_length", v.ContentLength),
					zap.Int64("response_size", v.ResponseSize),
					zap.String("latency", v.Latency.String()),
				)

				return nil
			},
		}))

		c.getHTTPHandler().RegisterRoutes(c.httpServer)
	}

	return c.httpServer
}

func (c *diContainer) getHTTPHandler() api.HTTPHandler {
	if c.httpHandler == nil {
		c.httpHandler = api.NewHandler(c.getPackCalculator(), c.getPackCache())
	}

	return c.httpHandler
}

func (c *diContainer) getPackCalculator() pack.Calculator {
	if c.packCalculator == nil {
		c.packCalculator = pack.NewCalculator(c.getConfig().OrderSizeLimit)
	}

	return c.packCalculator
}

func (c *diContainer) getPackCache() pack.Cache {
	if c.packCache == nil {
		cache, err := pack.NewInMemoryCache()
		if err != nil {
			zap.L().Fatal("failed to create pack cache", zap.Error(err))
		}

		c.packCache = cache
	}

	return c.packCache
}
