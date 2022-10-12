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
	role struct {
		Factory factory.Factory
	}
	Role interface {
		Route(g *echo.Group)
		Get(c echo.Context) error
	}
)

func NewRole(f factory.Factory) Role {
	return &role{f}
}

func (h *role) Route(g *echo.Group) {
	g.GET("", h.Get, middleware.Authentication)
}

// Get role
// @Summary Get role
// @Description Get role
// @Tags role
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query abstraction.Filter true "request query"
// @Param entity query model.RoleEntity false "entity query"
// @Success 200 {object} dto.RoleResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/roles [get]
func (h *role) Get(c echo.Context) error {
	filter := abstraction.NewFilterBuiler[model.RoleEntity](c, "roles")
	if err := c.Bind(filter.Payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	filter.Bind()

	result, pagination, err := h.Factory.Usecase.Role.Find(c.Request().Context(), *filter.Payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.CustomSuccessBuilder(200, result, "Get roles success", &pagination).Send(c)
}
