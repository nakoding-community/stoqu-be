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
	reminderStockHistory struct {
		Factory factory.Factory
	}
	ReminderStockHistory interface {
		Route(g *echo.Group)
		Get(c echo.Context) error
		GetByID(c echo.Context) error
		Update(c echo.Context) error
		Delete(c echo.Context) error
	}
)

func NewReminderStockHistory(f factory.Factory) ReminderStockHistory {
	return &reminderStockHistory{f}
}

func (h *reminderStockHistory) Route(g *echo.Group) {
	g.GET("", h.Get, middleware.Authentication)
	g.GET("/count-unread", h.GetCountUnead, middleware.Authentication)
	g.GET("/:id", h.GetByID, middleware.Authentication)
	g.PUT("/bulk-read", h.UpdateBulkRead, middleware.Authentication)
	g.PUT("/:id", h.Update, middleware.Authentication)
	g.DELETE("/:id", h.Delete, middleware.Authentication)
}

// Get reminderStockHistory
// @Summary Get reminderStockHistory
// @Description Get reminderStockHistory
// @Tags reminderStockHistory
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query abstraction.Filter true "request query"
// @Param entity query model.ReminderStockHistoryEntity false "entity query"
// @Success 200 {object} dto.ReminderStockHistoryResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/reminder-stock-histories [get]
func (h *reminderStockHistory) Get(c echo.Context) error {
	filter := abstraction.NewFilterBuiler[model.ReminderStockHistoryEntity](c, "reminder_stocks")
	if err := c.Bind(filter.Payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	filter.Bind()

	result, pagination, err := h.Factory.Usecase.ReminderStockHistory.Find(c.Request().Context(), *filter.Payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.CustomSuccessBuilder(200, result, "Get reminder stock success", &pagination).Send(c)
}

// Get count unread reminderStockHistory
// @Summary Get count unread reminderStockHistory
// @Description Get count unread reminderStockHistory
// @Tags reminderStockHistory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ReminderStockHistoryCountUnreadResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/reminder-stock-histories/count-unread [get]
func (h *reminderStockHistory) GetCountUnead(c echo.Context) error {
	ctx := c.Request().Context()
	result, err := h.Factory.Usecase.ReminderStockHistory.CountUnread(ctx)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}
	return res.SuccessResponse(result).Send(c)
}

// Get reminderStockHistory by id
// @Summary Get reminderStockHistory by id
// @Description Get reminderStockHistory by id
// @Tags reminderStockHistory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "id path"
// @Success 200 {object} dto.ReminderStockHistoryResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/reminder-stock-histories/{id} [get]
func (h *reminderStockHistory) GetByID(c echo.Context) error {
	ctx := c.Request().Context()

	payload := new(dto.ByIDRequest)
	if err := c.Bind(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		response := res.ErrorBuilder(res.Constant.Error.Validation, err)
		return response.Send(c)
	}

	result, err := h.Factory.Usecase.ReminderStockHistory.FindByID(ctx, *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}
	return res.SuccessResponse(result).Send(c)
}

// Update bulk read reminderStockHistory
// @Summary Update bulk read reminderStockHistory
// @Description Update bulk read reminderStockHistory
// @Tags reminderStockHistory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateReminderStockHistoryBulkReadRequest true "request body"
// @Success 200 {object} dto.ReminderStockHistoryBulkReadResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/reminder-stock-histories/bulk-read [put]
func (h *reminderStockHistory) UpdateBulkRead(c echo.Context) error {
	ctx := c.Request().Context()

	payload := new(dto.UpdateReminderStockHistoryBulkReadRequest)
	if err := c.Bind(&payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.Validation, err).Send(c)
	}

	result, err := h.Factory.Usecase.ReminderStockHistory.UpdateBulkRead(ctx, *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(result).Send(c)
}

// Update reminderStockHistory
// @Summary Update reminderStockHistory
// @Description Update reminderStockHistory
// @Tags reminderStockHistory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "id path"
// @Param request body dto.UpdateReminderStockHistoryRequest true "request body"
// @Success 200 {object} dto.ReminderStockHistoryResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/reminder-stock-histories/{id} [put]
func (h *reminderStockHistory) Update(c echo.Context) error {
	ctx := c.Request().Context()

	payload := new(dto.UpdateReminderStockHistoryRequest)
	if err := c.Bind(&payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.Validation, err).Send(c)
	}

	result, err := h.Factory.Usecase.ReminderStockHistory.Update(ctx, *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(result).Send(c)
}

// Delete reminderStockHistory
// @Summary Delete reminderStockHistory
// @Description Delete reminderStockHistory
// @Tags reminderStockHistory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "id path"
// @Success 200 {object} dto.ReminderStockHistoryResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/reminder-stock-histories/{id} [delete]
func (h *reminderStockHistory) Delete(c echo.Context) error {
	ctx := c.Request().Context()

	payload := new(dto.ByIDRequest)
	if err := c.Bind(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.Validation, err).Send(c)
	}

	result, err := h.Factory.Usecase.ReminderStockHistory.Delete(ctx, *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(result).Send(c)
}
