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
	productLookup struct {
		Factory factory.Factory
	}
	ProductLookup interface {
		Route(g *echo.Group)
		Get(c echo.Context) error
		GetByID(c echo.Context) error
		Create(c echo.Context) error
		Update(c echo.Context) error
		Delete(c echo.Context) error
	}
)

func NewProductLookup(f factory.Factory) ProductLookup {
	return &productLookup{f}
}

func (h *productLookup) Route(g *echo.Group) {
	g.GET("", h.Get, middleware.Authentication)
	g.GET("/:id", h.GetByID, middleware.Authentication)
	g.POST("", h.Create, middleware.Authentication)
	g.PUT("/:id", h.Update, middleware.Authentication)
	g.DELETE("/:id", h.Delete, middleware.Authentication)
}

// Get productLookup
// @Summary Get productLookup
// @Description Get productLookup
// @Tags productLookup
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query abstraction.Filter true "request query"
// @Param entity query model.ProductLookupEntity false "entity query"
// @Success 200 {object} dto.ProductLookupResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/product-lookups [get]
func (h *productLookup) Get(c echo.Context) error {
	filter := abstraction.NewFilterBuiler[model.ProductLookupEntity](c, "product_lookups")
	if err := c.Bind(filter.Payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	filter.Bind()

	result, pagination, err := h.Factory.Usecase.ProductLookup.Find(c.Request().Context(), *filter.Payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.CustomSuccessBuilder(200, result, "Get product lookups success", &pagination).Send(c)
}

// Get productLookup by id
// @Summary Get productLookup by id
// @Description Get productLookup by id
// @Tags productLookup
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "id path"
// @Success 200 {object} dto.ProductLookupResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/product-lookups/{id} [get]
func (h *productLookup) GetByID(c echo.Context) error {
	ctx := c.Request().Context()

	payload := new(dto.ByIDRequest)
	if err := c.Bind(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		response := res.ErrorBuilder(res.Constant.Error.Validation, err)
		return response.Send(c)
	}

	result, err := h.Factory.Usecase.ProductLookup.FindByID(ctx, *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}
	return res.SuccessResponse(result).Send(c)
}

// Create productLookup
// @Summary Create productLookup
// @Description Create productLookup
// @Tags productLookup
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateProductLookupRequest true "request body"
// @Success 200 {object} dto.ProductLookupResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/product-lookups [post]
func (h *productLookup) Create(c echo.Context) error {
	ctx := c.Request().Context()

	payload := new(dto.CreateProductLookupRequest)
	if err := c.Bind(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.Validation, err).Send(c)
	}

	result, err := h.Factory.Usecase.ProductLookup.Create(ctx, *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(result).Send(c)
}

// Update productLookup
// @Summary Update productLookup
// @Description Update productLookup
// @Tags productLookup
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "id path"
// @Param request body dto.UpdateProductLookupRequest true "request body"
// @Success 200 {object} dto.ProductLookupResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/product-lookups/{id} [put]
func (h *productLookup) Update(c echo.Context) error {
	ctx := c.Request().Context()

	payload := new(dto.UpdateProductLookupRequest)
	if err := c.Bind(&payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.Validation, err).Send(c)
	}

	result, err := h.Factory.Usecase.ProductLookup.Update(ctx, *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(result).Send(c)
}

// Delete productLookup
// @Summary Delete productLookup
// @Description Delete productLookup
// @Tags productLookup
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "id path"
// @Success 200 {object} dto.ProductLookupResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/product-lookups/{id} [delete]
func (h *productLookup) Delete(c echo.Context) error {
	ctx := c.Request().Context()

	payload := new(dto.ByIDRequest)
	if err := c.Bind(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.Validation, err).Send(c)
	}

	result, err := h.Factory.Usecase.ProductLookup.Delete(ctx, *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(result).Send(c)
}
