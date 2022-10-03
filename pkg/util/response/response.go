package response

import "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"

type Meta struct {
	Success bool                        `json:"success" default:"true"`
	Code    string                      `json:"code"`
	Message string                      `json:"message" default:"true"`
	Info    *abstraction.PaginationInfo `json:"info"`
}

type responseHelper struct {
	Error   errorHelper
	Success successHelper
}

var Constant responseHelper = responseHelper{
	Error:   errorConstant,
	Success: successConstant,
}
