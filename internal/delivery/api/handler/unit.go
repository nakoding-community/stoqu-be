package handler

import (
	"gitlab.com/stoqu/stoqu-be/internal/delivery/api/middleware"
	"gitlab.com/stoqu/stoqu-be/internal/factory"
	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"

	"github.com/labstack/echo/v4"
)

type (
	unit struct {
		Factory factory.Factory
	}
	Unit interface {
		Route(g *echo.Group)
		Get(c echo.Context) error
	}
)

func NewUnit(f factory.Factory) Unit {
	return &unit{f}
}

func (h *unit) Route(g *echo.Group) {
	g.GET("", h.Get, middleware.Authentication)
}

// Get unit
// @Summary Get unit
// @Description Get unit
// @Tags unit
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query abstraction.Filter true "request query"
// @Param entity query model.UnitEntity false "entity query"
// @Success 200 {object} dto.UnitResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/units [get]
func (h *unit) Get(c echo.Context) error {
	filter := abstraction.NewFilterBuiler[model.UnitEntity](c, "units")
	if err := c.Bind(filter.Payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	filter.Bind()

	result, pagination, err := h.Factory.Usecase.Unit.Find(c.Request().Context(), *filter.Payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.CustomSuccessBuilder(200, result, "Get units success", &pagination).Send(c)
}
