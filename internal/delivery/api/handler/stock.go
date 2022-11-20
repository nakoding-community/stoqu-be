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
	stock struct {
		Factory factory.Factory
	}
	Stock interface {
		Route(g *echo.Group)
		Get(c echo.Context) error
		GetByID(c echo.Context) error
		Transaction(c echo.Context) error
		Convertion(c echo.Context) error
		Movement(c echo.Context) error
	}
)

func NewStock(f factory.Factory) Stock {
	return &stock{f}
}

func (h *stock) Route(g *echo.Group) {
	g.GET("", h.Get, middleware.Authentication)
	g.GET("/:id", h.GetByID, middleware.Authentication)
	g.PUT("/transaction", h.Transaction, middleware.Authentication)
	g.PUT("/convertion", h.Convertion, middleware.Authentication)
	g.PUT("/movement", h.Movement, middleware.Authentication)
}

// Get stock
// @Summary Get stock
// @Description Get stock
// @Tags stock
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query abstraction.Filter true "request query"
// @Param entity query model.StockView false "entity query"
// @Success 200 {object} dto.StockResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/stocks [get]
func (h *stock) Get(c echo.Context) error {
	filter := abstraction.NewFilterBuiler[model.StockView](c, "stocks")
	if err := c.Bind(filter.Payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	filter.Bind()

	result, pagination, err := h.Factory.Usecase.Stock.Find(c.Request().Context(), *filter.Payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.CustomSuccessBuilder(200, result, "Get stocks success", &pagination).Send(c)
}

// Get stock history
// @Summary Get stock history
// @Description Get stock history
// @Tags stock history
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query abstraction.Filter true "request query"
// @Param entity query model.StockTrx false "entity query"
// @Success 200 {object} dto.StockResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/stocks/history [get]
func (h *stock) History(c echo.Context) error {
	filter := abstraction.NewFilterBuiler[model.StockTrxModel](c, "stock_trxs")
	if err := c.Bind(filter.Payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	filter.Bind()

	result, pagination, err := h.Factory.Usecase.Stock.History(c.Request().Context(), *filter.Payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.CustomSuccessBuilder(200, result, "Get stock histories success", &pagination).Send(c)
}

// Get stock by id
// @Summary Get stock by id
// @Description Get stock by id
// @Tags stock
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "id path"
// @Success 200 {object} dto.StockResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/stocks/{id} [get]
func (h *stock) GetByID(c echo.Context) error {
	ctx := c.Request().Context()

	payload := new(dto.ByIDRequest)
	if err := c.Bind(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		response := res.ErrorBuilder(res.Constant.Error.Validation, err)
		return response.Send(c)
	}

	result, err := h.Factory.Usecase.Stock.FindByID(ctx, *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}
	return res.SuccessResponse(result).Send(c)
}

// Transaction stock
// @Summary Transaction stock
// @Description Transaction stock
// @Tags stock
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.TransactionStockRequest true "request body"
// @Success 200 {object} dto.StockTransactionResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/stocks/transaction [put]
func (h *stock) Transaction(c echo.Context) error {
	ctx := c.Request().Context()

	payload := new(dto.TransactionStockRequest)
	if err := c.Bind(&payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.Validation, err).Send(c)
	}

	result, err := h.Factory.Usecase.Stock.Transaction(ctx, *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(result).Send(c)
}

// Convertion stock
// @Summary Convertion stock
// @Description Convertion stock
// @Tags stock
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ConvertionStockRequest true "request body"
// @Success 200 {object} dto.StockConvertionResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/stocks/convertion [put]
func (h *stock) Convertion(c echo.Context) error {
	ctx := c.Request().Context()

	payload := new(dto.ConvertionStockRequest)
	if err := c.Bind(&payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.Validation, err).Send(c)
	}

	result, err := h.Factory.Usecase.Stock.Convertion(ctx, *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(result).Send(c)
}

// Movement stock
// @Summary Movement stock
// @Description Movement stock
// @Tags stock
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.MovementStockRequest true "request body"
// @Success 200 {object} dto.StockMovementResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/stocks/convertion [put]
func (h *stock) Movement(c echo.Context) error {
	ctx := c.Request().Context()

	payload := new(dto.MovementStockRequest)
	if err := c.Bind(&payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.Validation, err).Send(c)
	}

	result, err := h.Factory.Usecase.Stock.Movement(ctx, *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(result).Send(c)
}
