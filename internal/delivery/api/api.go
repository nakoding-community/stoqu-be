package api

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/stoqu/stoqu-be/internal/delivery/api/handler"
	"gitlab.com/stoqu/stoqu-be/internal/factory"
)

func Init(e *echo.Echo, f factory.Factory) {
	// routes
	prefix := "api"
	handler.NewAuth(f).Route(e.Group(prefix + "/auth"))
	handler.NewUser(f).Route(e.Group(prefix + "/users"))
}
