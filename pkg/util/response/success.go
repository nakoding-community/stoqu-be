package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
)

type successHelper struct {
	OK Success
}

var successConstant successHelper = successHelper{
	OK: Success{
		Response: successResponse{
			Meta: Meta{
				Success: true,
				Message: "Request successfully proceed",
			},
			Data: nil,
		},
		HttpCode: http.StatusOK,
	},
}

type successResponse struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Success struct {
	Response successResponse `json:"response"`
	HttpCode int
}

func SuccessBuilder(res Success, data interface{}) *Success {
	res.Response.Data = data
	return &res
}

func CustomSuccessBuilder(code int, data interface{}, message string, info *abstraction.PaginationInfo) *Success {
	return &Success{
		Response: successResponse{
			Meta: Meta{
				Success: true,
				Message: message,
				Info:    info,
			},
			Data: data,
		},
		HttpCode: code,
	}
}

func SuccessResponse(data interface{}) *Success {
	return SuccessBuilder(Constant.Success.OK, data)
}

func (s *Success) Send(c echo.Context) error {
	return c.JSON(s.HttpCode, s.Response)
}
