package handler

import (
	"gitlab.com/stoqu/stoqu-be/internal/delivery/api/middleware"
	"gitlab.com/stoqu/stoqu-be/internal/factory"
	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	"gitlab.com/stoqu/stoqu-be/internal/model/dto"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"

	"github.com/labstack/echo/v4"
)

type (
	order struct {
		Factory factory.Factory
	}
	Order interface {
		Route(g *echo.Group)
		Get(c echo.Context) error
		GetByID(c echo.Context) error
		Upsert(c echo.Context) error
	}
)

func NewOrder(f factory.Factory) Order {
	return &order{f}
}

func (h *order) Route(g *echo.Group) {
	g.GET("", h.Get, middleware.Authentication)
	g.GET("/:id", h.GetByID, middleware.Authentication)
	g.PUT("/", h.Upsert, middleware.Authentication)
}

// Get order
// @Summary Get order
// @Description Get order
// @Tags order
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query abstraction.Filter true "request query"
// @Param entity query model.OrderView false "entity query"
// @Success 200 {object} dto.OrderViewResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/orders [get]
func (h *order) Get(c echo.Context) error {
	filter := abstraction.NewFilterBuiler[model.OrderView](c, "orders")
	if err := c.Bind(filter.Payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	filter.Bind()

	result, pagination, err := h.Factory.Usecase.Order.Find(c.Request().Context(), *filter.Payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.CustomSuccessBuilder(200, result, "Get orders success", &pagination).Send(c)
}

// Get order by id
// @Summary Get order by id
// @Description Get order by id
// @Tags order
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "id path"
// @Success 200 {object} dto.OrderViewResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/orders/{id} [get]
func (h *order) GetByID(c echo.Context) error {
	ctx := c.Request().Context()

	payload := new(dto.ByIDRequest)
	if err := c.Bind(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		response := res.ErrorBuilder(res.Constant.Error.Validation, err)
		return response.Send(c)
	}

	result, err := h.Factory.Usecase.Order.FindByID(ctx, *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}
	return res.SuccessResponse(result).Send(c)
}

// Upsert order
// @Summary Upsert order
// @Description Upsert order
// @Tags order
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpsertOrderRequest true "request body"
// @Success 200 {object} dto.OrderUpsertResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/orders [put]
func (h *order) Upsert(c echo.Context) error {
	ctx := c.Request().Context()

	payload := new(dto.UpsertOrderRequest)
	if err := c.Bind(&payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.Validation, err).Send(c)
	}

	result, err := h.Factory.Usecase.Order.Upsert(ctx, *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(result).Send(c)
}
