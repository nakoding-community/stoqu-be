package api

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/stoqu/stoqu-be/internal/delivery/api/handler"
	"gitlab.com/stoqu/stoqu-be/internal/factory"
)

func Init(e *echo.Echo, f factory.Factory) {
	// routes
	prefix := "api"

	handler.NewRole(f).Route(e.Group(prefix + "/roles"))
	handler.NewUser(f).Route(e.Group(prefix + "/users"))
	handler.NewAuth(f).Route(e.Group(prefix + "/auth"))

	handler.NewUnit(f).Route(e.Group(prefix + "/units"))
	handler.NewConvertionUnit(f).Route(e.Group(prefix + "/convertion-units"))
	handler.NewCurrency(f).Route(e.Group(prefix + "/currencies"))
	handler.NewCurrency(f).Route(e.Group(prefix + "/reminder-stocks"))

	handler.NewPacket(f).Route(e.Group(prefix + "/packets"))
	handler.NewBrand(f).Route(e.Group(prefix + "/brands"))
	handler.NewVariant(f).Route(e.Group(prefix + "/variants"))
}
