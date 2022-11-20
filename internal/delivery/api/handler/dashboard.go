package handler

import (
	"gitlab.com/stoqu/stoqu-be/internal/delivery/api/middleware"
	"gitlab.com/stoqu/stoqu-be/internal/factory"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"

	"github.com/labstack/echo/v4"
)

type (
	dashboard struct {
		Factory factory.Factory
	}
	Dashboard interface {
		Route(g *echo.Group)
		Count(c echo.Context) error
	}
)

func NewDashboard(f factory.Factory) Dashboard {
	return &dashboard{f}
}

func (h *dashboard) Route(g *echo.Group) {
	g.GET("/count", h.Count, middleware.Authentication)
}

// Count dashboard
// @Summary Count dashboard
// @Description Count dashboard
// @Tags dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.DashboardResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/dashboards/count [get]
func (h *dashboard) Count(c echo.Context) error {
	result, err := h.Factory.Usecase.Dashboard.Count(c.Request().Context())
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(result).Send(c)
}
