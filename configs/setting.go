package configs

import (
	"context"
	"encoding/json"
	"fmt"
	"go-oms/configs/store"
	"go-oms/internal/middlewares"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

type Setting struct {
	App    *fiber.App
	Logger *zap.SugaredLogger
	PG     *store.PostgresqlStore
}

func NewApp(logger *zap.SugaredLogger) *Setting {
	return &Setting{Logger: logger}
}

func (c *Setting) SetApp(ctx context.Context) error {
	c.Logger.Named("backend-chellenge")

	c.App = fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
	})

	if err := SetEnv(ctx); err != nil {
		return err
	}

	cfg := cors.Config{
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
	}

	c.App.Use(recover.New(recover.Config{EnableStackTrace: true}))
	c.App.Use(helmet.New())
	c.App.Use(cors.New(cfg))
	c.App.Use(middlewares.LoggerMiddleware(c.Logger))
	c.App.Use(limiter.New(limiter.Config{
		Max:        1000,
		Expiration: time.Second,
	}))
	c.App.Use(logger.New(logger.Config{
		Format:     "${blue}${time} ${yellow}${status} - ${red}${latency} ${cyan}${method} ${path} ${green} ${ip} ${ua} ${reset}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Asia/Bangkok",
		Output:     os.Stdout,
	}))

	postgres, err := store.ConnectPostgresql()
	if err != nil {
		return err
	}

	c.PG = postgres
	return nil
}

func (c *Setting) RunApp(ctx context.Context) <-chan error {
	errChan := make(chan error, 1)

	go func() {
		err := c.App.Listen(fmt.Sprintf(":%v", App.Port))
		if err != nil {
			errChan <- fmt.Errorf("fiber listen error: %w", err)
		}
	}()

	return errChan
}

func (c *Setting) StopApp(ctx context.Context) error {
	if err := c.App.Shutdown(); err != nil {
		return fmt.Errorf("error app failed to stop : %s", err)
	}

	// ปิด PostgreSQL connection pool
	if c.PG != nil && c.PG.Store != nil {
		sqlDB, err := c.PG.Store.DB()
		if err != nil {
			return fmt.Errorf("failed to get sql.DB from GORM: %w", err)
		}
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
	}

	return nil
}
