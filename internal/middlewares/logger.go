package middlewares

import (
	"context"
	"time"

	"go-oms/internal/domain/entities"
	pkglogger "go-oms/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func LoggerMiddleware(logger *zap.SugaredLogger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := uuid.New().String()
		ctx := c.UserContext()
		logger := logger.With("request_id", requestID)
		ctx = pkglogger.WithLogger(ctx, logger)
		ctx = context.WithValue(ctx, entities.RequestId, requestID)
		c.SetUserContext(ctx)

		start := time.Now()
		err := c.Next()
		stop := time.Now()
		latency := stop.Sub(start)

		logger.Infow("HTTP Request",
			"method", c.Method(),
			"path", c.OriginalURL(),
			"status", c.Response().StatusCode(),
			"latency", latency.String(),
			"ip", c.IP(),
			"user_agent", c.Get("User-Agent"),
		)

		return err
	}
}
