package middleware

import (
	"strconv"

	"gitlab.com/stoqu/stoqu-be/pkg/util/response"

	"github.com/labstack/echo/v4"
)

func Mock(reqDto interface{}, resDto interface{}) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			mockS := c.Request().Header.Get("Mock")
			if mockS != "" {
				//checking
				if v, err := strconv.ParseBool(mockS); err == nil && v {
					if err := c.Bind(reqDto); err != nil {
						return response.ErrorBuilder(response.Constant.Error.BadRequest, err).Send(c)
					}

					return response.SuccessResponse(resDto).Send(c)
				}
			}

			return next(c)
		}
	}
}
