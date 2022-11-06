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
	currency struct {
		Factory factory.Factory
	}
	Currency interface {
		Route(g *echo.Group)
		Get(c echo.Context) error
		GetByID(c echo.Context) error
		Create(c echo.Context) error
		Update(c echo.Context) error
		Delete(c echo.Context) error
		Convert(c echo.Context) error
	}
)

func NewCurrency(f factory.Factory) Currency {
	return &currency{f}
}

func (h *currency) Route(g *echo.Group) {
	g.GET("", h.Get, middleware.Authentication)
	g.GET("/:id", h.GetByID, middleware.Authentication)
	g.POST("", h.Create, middleware.Authentication)
	g.PUT("/:id", h.Update, middleware.Authentication)
	g.DELETE("/:id", h.Delete, middleware.Authentication)
	g.POST("/convert", h.Convert, middleware.Authentication)
}

// Get currency
// @Summary Get currency
// @Description Get currency
// @Tags currency
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query abstraction.Filter true "request query"
// @Param entity query model.CurrencyEntity false "entity query"
// @Success 200 {object} dto.CurrencyResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/currencies [get]
func (h *currency) Get(c echo.Context) error {
	filter := abstraction.NewFilterBuiler[model.CurrencyEntity](c, "currencies")
	if err := c.Bind(filter.Payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	filter.Bind()

	result, pagination, err := h.Factory.Usecase.Currency.Find(c.Request().Context(), *filter.Payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.CustomSuccessBuilder(200, result, "Get currencies success", &pagination).Send(c)
}

// Get currency by id
// @Summary Get currency by id
// @Description Get currency by id
// @Tags currency
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "id path"
// @Success 200 {object} dto.CurrencyResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/currencies/{id} [get]
func (h *currency) GetByID(c echo.Context) error {
	ctx := c.Request().Context()

	payload := new(dto.ByIDRequest)
	if err := c.Bind(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		response := res.ErrorBuilder(res.Constant.Error.Validation, err)
		return response.Send(c)
	}

	result, err := h.Factory.Usecase.Currency.FindByID(ctx, *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}
	return res.SuccessResponse(result).Send(c)
}

// Create currency
// @Summary Create currency
// @Description Create currency
// @Tags currency
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateCurrencyRequest true "request body"
// @Success 200 {object} dto.CurrencyResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/currencies [post]
func (h *currency) Create(c echo.Context) error {
	ctx := c.Request().Context()

	payload := new(dto.CreateCurrencyRequest)
	if err := c.Bind(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.Validation, err).Send(c)
	}

	result, err := h.Factory.Usecase.Currency.Create(ctx, *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(result).Send(c)
}

// Update currency
// @Summary Update currency
// @Description Update currency
// @Tags currency
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "id path"
// @Param request body dto.UpdateCurrencyRequest true "request body"
// @Success 200 {object} dto.CurrencyResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/currencies/{id} [put]
func (h *currency) Update(c echo.Context) error {
	ctx := c.Request().Context()

	payload := new(dto.UpdateCurrencyRequest)
	if err := c.Bind(&payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.Validation, err).Send(c)
	}

	result, err := h.Factory.Usecase.Currency.Update(ctx, *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(result).Send(c)
}

// Delete currency
// @Summary Delete currency
// @Description Delete currency
// @Tags currency
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "id path"
// @Success 200 {object} dto.CurrencyResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/currencies/{id} [delete]
func (h *currency) Delete(c echo.Context) error {
	ctx := c.Request().Context()

	payload := new(dto.ByIDRequest)
	if err := c.Bind(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.Validation, err).Send(c)
	}

	result, err := h.Factory.Usecase.Currency.Delete(ctx, *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(result).Send(c)
}

// Convert currency
// @Summary Convert currency
// @Description Convert currency
// @Tags currency
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ConvertCurrencyRequest true "request body"
// @Success 200 {object} dto.ConvertCurrencyResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/currencies/convert [post]
func (h *currency) Convert(c echo.Context) error {
	ctx := c.Request().Context()

	payload := new(dto.ConvertCurrencyRequest)
	if err := c.Bind(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.Validation, err).Send(c)
	}

	result, err := h.Factory.Usecase.Currency.Convert(ctx, *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(result).Send(c)
}
