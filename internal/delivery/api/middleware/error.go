package middleware

import (
	"net/http"

	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {
	var errCustom *res.Error

	report, ok := err.(*echo.HTTPError)
	if !ok {
		report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	switch report.Code {
	case http.StatusNotFound:
		errCustom = res.ErrorBuilder(res.Constant.Error.RouteNotFound, err)
	default:
		errCustom = res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}

	res.ErrorResponse(errCustom).Send(c)
}
