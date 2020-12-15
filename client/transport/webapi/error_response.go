package webapi

import (
	"encoding/json"
)

type apiError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func newApiError(err error, code int) []byte {
	errResp := apiError{
		Error: err.Error(),
		Code:  code,
	}
	resp, _ := json.Marshal(errResp)
	return resp
}
