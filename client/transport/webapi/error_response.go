package webapi

import (
	"encoding/json"
)

// apiError a set of fields to show error
type apiError struct {
	// Error - error message from api
	Error string `json:"error"`
	// Code - http code
	Code int `json:"code"`
}

func newApiError(err error, code int) []byte {
	errResp := apiError{
		Error: err.Error(),
		Code:  code,
	}
	resp, _ := json.Marshal(errResp)
	return resp
}
