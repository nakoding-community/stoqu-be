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
	reminderStock struct {
		Factory factory.Factory
	}
	ReminderStock interface {
		Route(g *echo.Group)
		Get(c echo.Context) error
		GetByID(c echo.Context) error
		Create(c echo.Context) error
		Update(c echo.Context) error
		Delete(c echo.Context) error
	}
)

func NewReminderStock(f factory.Factory) ReminderStock {
	return &reminderStock{f}
}

func (h *reminderStock) Route(g *echo.Group) {
	g.GET("", h.Get, middleware.Authentication)
	g.GET("/:id", h.GetByID, middleware.Authentication)
	g.POST("", h.Create, middleware.Authentication)
	g.PUT("/:id", h.Update, middleware.Authentication)
	g.DELETE("/:id", h.Delete, middleware.Authentication)
}

// Get reminderStock
// @Summary Get reminderStock
// @Description Get reminderStock
// @Tags reminderStock
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query abstraction.Filter true "request query"
// @Param entity query model.ReminderStockEntity false "entity query"
// @Success 200 {object} dto.ReminderStockResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/reminder_stocks [get]
func (h *reminderStock) Get(c echo.Context) error {
	filter := abstraction.NewFilterBuiler[model.ReminderStockEntity](c, "reminder_stocks")
	if err := c.Bind(filter.Payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	filter.Bind()

	result, pagination, err := h.Factory.Usecase.ReminderStock.Find(c.Request().Context(), *filter.Payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.CustomSuccessBuilder(200, result, "Get reminder stock success", &pagination).Send(c)
}

// Get reminderStock by id
// @Summary Get reminderStock by id
// @Description Get reminderStock by id
// @Tags reminderStock
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "id path"
// @Success 200 {object} dto.ReminderStockResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/reminder-stocks/{id} [get]
func (h *reminderStock) GetByID(c echo.Context) error {
	ctx := c.Request().Context()

	payload := new(dto.ByIDRequest)
	if err := c.Bind(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		response := res.ErrorBuilder(res.Constant.Error.Validation, err)
		return response.Send(c)
	}

	result, err := h.Factory.Usecase.ReminderStock.FindByID(ctx, *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}
	return res.SuccessResponse(result).Send(c)
}

// Create reminderStock
// @Summary Create reminderStock
// @Description Create reminderStock
// @Tags reminderStock
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateReminderStockRequest true "request body"
// @Success 200 {object} dto.ReminderStockResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/reminder-stocks [post]
func (h *reminderStock) Create(c echo.Context) error {
	ctx := c.Request().Context()

	payload := new(dto.CreateReminderStockRequest)
	if err := c.Bind(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.Validation, err).Send(c)
	}

	result, err := h.Factory.Usecase.ReminderStock.Create(ctx, *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(result).Send(c)
}

// Update reminderStock
// @Summary Update reminderStock
// @Description Update reminderStock
// @Tags reminderStock
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "id path"
// @Param request body dto.UpdateReminderStockRequest true "request body"
// @Success 200 {object} dto.ReminderStockResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/reminder-stocks/{id} [put]
func (h *reminderStock) Update(c echo.Context) error {
	ctx := c.Request().Context()

	payload := new(dto.UpdateReminderStockRequest)
	if err := c.Bind(&payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.Validation, err).Send(c)
	}

	result, err := h.Factory.Usecase.ReminderStock.Update(ctx, *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(result).Send(c)
}

// Delete reminderStock
// @Summary Delete reminderStock
// @Description Delete reminderStock
// @Tags reminderStock
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "id path"
// @Success 200 {object} dto.ReminderStockResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/reminder-stocks/{id} [delete]
func (h *reminderStock) Delete(c echo.Context) error {
	ctx := c.Request().Context()

	payload := new(dto.ByIDRequest)
	if err := c.Bind(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.Validation, err).Send(c)
	}

	result, err := h.Factory.Usecase.ReminderStock.Delete(ctx, *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(result).Send(c)
}
