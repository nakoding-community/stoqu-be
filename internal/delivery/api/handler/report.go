package handler

import (
	"gitlab.com/stoqu/stoqu-be/internal/delivery/api/middleware"
	"gitlab.com/stoqu/stoqu-be/internal/factory"
	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	"gitlab.com/stoqu/stoqu-be/internal/model/dto"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"

	"github.com/labstack/echo/v4"
)

type (
	report struct {
		Factory factory.Factory
	}
	Report interface {
		Route(g *echo.Group)
		GetOrder(c echo.Context) error
	}
)

func NewReport(f factory.Factory) Report {
	return &report{f}
}

func (h *report) Route(g *echo.Group) {
	g.GET("/orders", h.GetOrder, middleware.Authentication)
}

// Get reportOrders
// @Summary Get report orders
// @Description Get report orders
// @Tags report
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query abstraction.Filter true "request query"
// @Param entity query dto.OrderReportQuery false "entity query"
// @Success 200 {object} dto.OrderReportResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /api/reports/orders [get]
func (h *report) GetOrder(c echo.Context) error {
	filter := abstraction.NewFilterBuiler[dto.OrderReportQuery](c, "order_trxs")
	if err := c.Bind(filter.Payload); err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err).Send(c)
	}
	filter.Bind()

	result, pagination, err := h.Factory.Usecase.Report.FindOrder(c.Request().Context(), *filter.Payload, dto.OrderReportQuery{})
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.CustomSuccessBuilder(200, result, "Get order reports success", &pagination).Send(c)
}
