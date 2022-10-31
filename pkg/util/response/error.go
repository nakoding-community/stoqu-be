package response

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	resConstant "gitlab.com/stoqu/stoqu-be/pkg/util/response/constant"
)

type errorResponse struct {
	Meta  Meta   `json:"meta"`
	Error string `json:"error"`
}

type Error struct {
	Response     errorResponse `json:"response"`
	HttpCode     int
	ErrorMessage error
}

type errorHelper struct {
	Duplicate           Error
	NotFound            Error
	RouteNotFound       Error
	UnprocessableEntity Error
	Unauthorized        Error
	BadRequest          Error
	Validation          Error
	InternalServerError Error
}

var errorConstant errorHelper = errorHelper{
	Duplicate: Error{
		Response: errorResponse{
			Meta: Meta{
				Success: false,
				Code:    resConstant.E_DUPLICATE,
				Message: resConstant.CODE_TO_MESSAGE[resConstant.E_DUPLICATE],
			},
			Error: resConstant.CODE_TO_MESSAGE[resConstant.E_DUPLICATE],
		},
		HttpCode: http.StatusConflict,
	},
	NotFound: Error{
		Response: errorResponse{
			Meta: Meta{
				Success: false,
				Message: resConstant.CODE_TO_MESSAGE[resConstant.E_DATA_NOTFOUND],
				Code:    resConstant.E_DATA_NOTFOUND,
			},
			Error: resConstant.CODE_TO_MESSAGE[resConstant.E_DATA_NOTFOUND],
		},
		HttpCode: http.StatusNotFound,
	},
	RouteNotFound: Error{
		Response: errorResponse{
			Meta: Meta{
				Success: false,
				Message: "Route not found",
				Code:    resConstant.E_ROUTE_NOTFOUND,
			},
			Error: resConstant.CODE_TO_MESSAGE[resConstant.E_ROUTE_NOTFOUND],
		},
		HttpCode: http.StatusNotFound,
	},
	UnprocessableEntity: Error{
		Response: errorResponse{
			Meta: Meta{
				Success: false,
				Message: resConstant.CODE_TO_MESSAGE[resConstant.E_UNPROCESSABL_ENTITY],
				Code:    resConstant.E_UNPROCESSABL_ENTITY,
			},
			Error: resConstant.CODE_TO_MESSAGE[resConstant.E_UNPROCESSABL_ENTITY],
		},
		HttpCode: http.StatusUnprocessableEntity,
	},
	Unauthorized: Error{
		Response: errorResponse{
			Meta: Meta{
				Success: false,
				Message: resConstant.CODE_TO_MESSAGE[resConstant.E_UNAUTHORIZED],
				Code:    resConstant.E_UNAUTHORIZED,
			},
			Error: resConstant.CODE_TO_MESSAGE[resConstant.E_UNAUTHORIZED],
		},
		HttpCode: http.StatusUnauthorized,
	},
	BadRequest: Error{
		Response: errorResponse{
			Meta: Meta{
				Success: false,
				Message: resConstant.CODE_TO_MESSAGE[resConstant.E_BAD_REQUEST],
				Code:    resConstant.E_BAD_REQUEST,
			},
			Error: resConstant.CODE_TO_MESSAGE[resConstant.E_BAD_REQUEST],
		},
		HttpCode: http.StatusBadRequest,
	},
	Validation: Error{
		Response: errorResponse{
			Meta: Meta{
				Success: false,
				Message: resConstant.CODE_TO_MESSAGE[resConstant.E_VALIDATION],
				Code:    resConstant.E_VALIDATION,
			},
			Error: resConstant.CODE_TO_MESSAGE[resConstant.E_VALIDATION],
		},
		HttpCode: http.StatusBadRequest,
	},
	InternalServerError: Error{
		Response: errorResponse{
			Meta: Meta{
				Success: false,
				Message: resConstant.CODE_TO_MESSAGE[resConstant.E_SERVER_ERROR],
				Code:    resConstant.E_SERVER_ERROR,
			},
			Error: resConstant.CODE_TO_MESSAGE[resConstant.E_SERVER_ERROR],
		},
		HttpCode: http.StatusInternalServerError,
	},
}

func ErrorBuilder(res Error, err error, customMessage ...string) *Error {
	res.ErrorMessage = err
	res.Response.Error = err.Error()

	if strings.Contains(strings.Join([]string{
		resConstant.CODE_TO_MESSAGE[resConstant.E_VALIDATION],
		resConstant.CODE_TO_MESSAGE[resConstant.E_BAD_REQUEST],
		resConstant.CODE_TO_MESSAGE[resConstant.E_DUPLICATE],
	}, ","), res.Response.Error) {
		res.Response.Meta.Message = err.Error()
	}

	if len(customMessage) > 0 {
		res.Response.Meta.Message = ""
		for i := range customMessage {
			res.Response.Meta.Message += customMessage[i]
		}
	}

	return &res
}

func CustomErrorBuilder(httpCode int, code string, err error, message string) *Error {
	return &Error{
		Response: errorResponse{
			Meta: Meta{
				Success: false,
				Message: message,
				Code:    code,
			},
			Error: err.Error(),
		},
		HttpCode: httpCode,
	}
}

func ErrorResponse(err error) *Error {
	re, ok := err.(*Error)
	if ok {
		return re
	} else {
		return ErrorBuilder(Constant.Error.InternalServerError, err)
	}
}

func (e *Error) Error() string {
	return e.ErrorMessage.Error()
}

func (e *Error) ParseToError() error {
	return e
}

func (e *Error) ErrorCode() string {
	return e.Response.Meta.Code
}

func (e *Error) Send(c echo.Context) error {
	if e.Response.Meta.Code != "" && e.Response.Meta.Message == "" {
		e.Response.Meta.Message = resConstant.CODE_TO_MESSAGE[e.Response.Meta.Code]
	} else if e.Response.Meta.Code == "" && e.Response.Meta.Message != "" {
		e.Response.Meta.Code = resConstant.MESSAGE_TO_CODE(e.Response.Meta.Message)
	}

	return c.JSON(e.HttpCode, e.Response)
}
