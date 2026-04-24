package main

import (
	"net/http"
	"os"
	"soliant-mock-api/modules/auth"

	"soliant-mock-api/modules/availability"
	"soliant-mock-api/modules/home"
	"soliant-mock-api/modules/location"
	"soliant-mock-api/modules/notifications"
	"soliant-mock-api/modules/opportunities"
	"soliant-mock-api/modules/pastshifts"
	"soliant-mock-api/modules/permissions"
	"soliant-mock-api/modules/profile"
	"soliant-mock-api/modules/shifts"
	"soliant-mock-api/modules/unavailableconfig"
	"soliant-mock-api/shared/httpmod"
	"soliant-mock-api/shared/logger"
	"soliant-mock-api/shared/redis"

	"github.com/MadAppGang/httplog"
	"github.com/MadAppGang/httplog/echolog"
	"github.com/labstack/echo/v4"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	if os.Getenv("ENV") == "PROD" {
		logger.Init("info")
	} else {
		logger.Init("debug")

		err := godotenv.Load()
		if err != nil {
			logger.Error(err.Error())
		}
	}

	e := echo.New()

	if os.Getenv("ENV") != "PROD" {
		e.Pre(echolog.LoggerWithConfig(httplog.LoggerConfig{
			RouterName:  "DEFAULT",
			CaptureBody: true,
		}))
	}

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-User-UUID"},
	}))

	cache, err := redis.NewClient("soliant-mock-api:")
	if err != nil {
		logger.Fatal("Failed to connect to redis: ", err)
	}

	shiftsMod := shifts.New(cache)

	// Register every feature module here. Adding a new CRUD resource is just
	// a matter of creating a package under modules/<name> (model.go, seed.go,
	// module.go) and appending it to this slice.
	modules := []httpmod.Module{
		shiftsMod,
		notifications.New(cache),
		availability.New(cache),
		location.New(cache),
		unavailableconfig.New(cache),
		home.New(cache),
		opportunities.New(cache),
		profile.New(cache),
		permissions.New(cache),
		auth.New(cache),
		pastshifts.New(shiftsMod.Repo),
	}

	for _, m := range modules {
		if err := m.SeedIfEmpty(); err != nil {
			logger.Fatal("Failed to seed module ", m.Name(), ": ", err)
		}
		m.Register(e)
		logger.Info("module registered: ", m.Name())
	}

	// Shared / infrastructure endpoints.
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "OK",
		})
	})

	e.POST("/logger", func(c echo.Context) error {
		var req struct {
			Level string `json:"level" binding:"required"`
		}

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		}

		logger.SetLevel(req.Level)

		return c.JSON(http.StatusOK, map[string]interface{}{"message": "OK"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Fatal(e.Start(port))
}
