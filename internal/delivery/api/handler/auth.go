package handler

import (
	"gitlab.com/stoqu/stoqu-be/internal/factory"
	"gitlab.com/stoqu/stoqu-be/internal/model/dto"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"

	"github.com/labstack/echo/v4"
)

type (
	auth struct {
		Factory factory.Factory
	}
	Auth interface {
		Route(g *echo.Group)
		Login(c echo.Context) error
		Register(c echo.Context) error
	}
)

func NewAuth(f factory.Factory) *auth {
	return &auth{f}
}

func (h *auth) Route(g *echo.Group) {
	g.POST("/login", h.Login)
	g.POST("/register", h.Register)
}

// Login godoc
// @Summary Login user
// @Description Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginAuthRequest true "request body"
// @Success 200 {object} dto.AuthLoginResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/auth/login [post]
func (h *auth) Login(c echo.Context) error {
	payload := new(dto.LoginAuthRequest)
	if err := c.Bind(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.Validation, err).Send(c)
	}

	data, err := h.Factory.Usecase.Auth.Login(c.Request().Context(), *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(data).Send(c)
}

// Register
// @Summary Register user
// @Description Register user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.RegisterAuthRequest true "request body"
// @Success 200 {object} dto.AuthLoginResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/auth/register [post]
func (h *auth) Register(c echo.Context) error {
	payload := new(dto.RegisterAuthRequest)
	if err := c.Bind(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.Validation, err).Send(c)
	}

	data, err := h.Factory.Usecase.Auth.Register(c.Request().Context(), *payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(data).Send(c)
}
