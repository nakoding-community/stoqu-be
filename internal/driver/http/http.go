package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gitlab.com/stoqu/stoqu-be/docs"
	"gitlab.com/stoqu/stoqu-be/internal/config"
	api "gitlab.com/stoqu/stoqu-be/internal/delivery/api"
	"gitlab.com/stoqu/stoqu-be/internal/delivery/api/middleware"
	"gitlab.com/stoqu/stoqu-be/internal/factory"
	"gitlab.com/stoqu/stoqu-be/pkg/util/gracefull"
)

func Init(cfg *config.Configuration, f factory.Factory) (gracefull.ProcessStarter, gracefull.ProcessStopper) {
	var (
		APP        = cfg.App.Name
		VERSION    = cfg.App.Version
		DOC_HOST   = cfg.Swagger.SwaggerHost
		DOC_SCHEME = cfg.Swagger.SwaggerScheme
	)
	// echo
	e := echo.New()

	// index
	e.GET("/", func(c echo.Context) error {
		message := fmt.Sprintf("Welcome to %s version %s", APP, VERSION)
		return c.String(http.StatusOK, message)
	})

	// doc
	docs.SwaggerInfo.Title = APP
	docs.SwaggerInfo.Version = VERSION
	docs.SwaggerInfo.Host = DOC_HOST
	docs.SwaggerInfo.Schemes = []string{DOC_SCHEME}
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// delivery
	middleware.Init(e)
	api.Init(e, f)
	// ws.Init(e, f)

	return func() error {
			return e.Start(":" + cfg.App.Port)
		}, func(ctx context.Context) error {
			return e.Shutdown(ctx)
		}
}
