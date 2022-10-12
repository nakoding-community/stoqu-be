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
	convertionUnit struct {
		Factory factory.Factory
	}
	ConvertionUnit interface {
		Route(g *echo.Group)
		Get(c echo.Context) error
		GetByID(c echo.Context) error
		Create(c echo.Context) error
		Update(c echo.Context) error
		Delete(c echo.Context) error
	}
)

func NewConvertionUnit(f factory.Factory) ConvertionUnit {
	return &convertionUnit{f}
}

func (h *convertionUnit) Route(g *echo.Group) {
	g.GET("", h.Get, middleware.Authentication)
	g.GET("/:id", h.GetByID, middleware.Authentication)
	g.POST("", h.Create, middleware.Authentication)
	g.PUT("/:id", h.Update, middleware.Authentication)
	g.DELETE("/:id", h.Delete, middleware.Authentication)
}

// Get convertionUnit
// @Summary Get convertionUnit
// @Description Get convertionUnit
// @Tags convertionUnit
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query abstraction.Filter true "request query"
// @Param entity query model.ConvertionUnitEntity false "entity query"
// @Success 200 {object} dto.ConvertionUnitResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/convertion-units [get]
func (h *convertionUnit) Get(c echo.Context) error {
	filter := abstraction.NewFilterBuiler[model.ConvertionUnitEntity](c, "convertion_units")
	if err := c.Bind(filter.Payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	filter.Bind()

	result, pagination, err := h.Factory.Usecase.ConvertionUnit.Find(c.Request().Context(), *filter.Payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.CustomSuccessBuilder(200, result, "Get convertion units success", &pagination).Send(c)
}

// Get convertionUnit by id
// @Summary Get convertionUnit by id
// @Description Get convertionUnit by id
// @Tags convertionUnit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "id path"
// @Success 200 {object} dto.ConvertionUnitResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/convertion-units/{id} [get]
func (h *convertionUnit) GetByID(c echo.Context) error {
	ctx := c.Request().Context()

	payload := new(dto.ByIDRequest)
	if err := c.Bind(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		response := res.ErrorBuilder(res.Constant.Error.Validation, err)
		return response.Send(c)
	}

	result, err := h.Factory.Usecase.ConvertionUnit.FindByID(ctx, *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}
	return res.SuccessResponse(result).Send(c)
}

// Create convertionUnit
// @Summary Create convertionUnit
// @Description Create convertionUnit
// @Tags convertionUnit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateConvertionUnitRequest true "request body"
// @Success 200 {object} dto.ConvertionUnitResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/convertion-units [post]
func (h *convertionUnit) Create(c echo.Context) error {
	ctx := c.Request().Context()

	payload := new(dto.CreateConvertionUnitRequest)
	if err := c.Bind(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.Validation, err).Send(c)
	}

	result, err := h.Factory.Usecase.ConvertionUnit.Create(ctx, *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(result).Send(c)
}

// Update convertionUnit
// @Summary Update convertionUnit
// @Description Update convertionUnit
// @Tags convertionUnit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "id path"
// @Param request body dto.UpdateConvertionUnitRequest true "request body"
// @Success 200 {object} dto.ConvertionUnitResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/convertion-units/{id} [put]
func (h *convertionUnit) Update(c echo.Context) error {
	ctx := c.Request().Context()

	payload := new(dto.UpdateConvertionUnitRequest)
	if err := c.Bind(&payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.Validation, err).Send(c)
	}

	result, err := h.Factory.Usecase.ConvertionUnit.Update(ctx, *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(result).Send(c)
}

// Delete convertionUnit
// @Summary Delete convertionUnit
// @Description Delete convertionUnit
// @Tags convertionUnit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "id path"
// @Success 200 {object} dto.ConvertionUnitResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/convertion-units/{id} [delete]
func (h *convertionUnit) Delete(c echo.Context) error {
	ctx := c.Request().Context()

	payload := new(dto.ByIDRequest)
	if err := c.Bind(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.Validation, err).Send(c)
	}

	result, err := h.Factory.Usecase.ConvertionUnit.Delete(ctx, *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(result).Send(c)
}
