package constant

const (
	//basic
	E_UNAUTHORIZED        = "E-4001"
	E_DUPLICATE           = "E-4002"
	E_BAD_REQUEST         = "E-4003"
	E_VALIDATION          = "E-4004"
	E_ROUTE_NOTFOUND      = "E-4005"
	E_DATA_NOTFOUND       = "E-4006"
	E_UNPROCESSABL_ENTITY = "E-4007"
	E_INVALID_CREDENTIALS = "E-4008"

	E_SERVER_ERROR = "E-5001"
)

var (
	CODE_TO_MESSAGE = map[string]string{
		E_UNAUTHORIZED:        "unauthorized, please login",
		E_DUPLICATE:           "created value already exists",
		E_BAD_REQUEST:         "bad request, please check payload",
		E_VALIDATION:          "invalid request parameters",
		E_ROUTE_NOTFOUND:      "route not found",
		E_DATA_NOTFOUND:       "data not found",
		E_SERVER_ERROR:        "internal server error",
		E_UNPROCESSABL_ENTITY: "unprocessable entity",
		E_INVALID_CREDENTIALS: "invalid credentials",
	}

	MESSAGE_TO_CODE = func(message string) string {
		for c, v := range CODE_TO_MESSAGE {
			if v == message {
				return c
			}
		}
		return E_SERVER_ERROR
	}
)
