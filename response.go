package v2rayaguard

import "strings"

type Response struct {
	Code string                 `json:"code"`
	Data map[string]interface{} `json:"data"`
}

func (r Response) isSuccess() bool {
	return strings.ToUpper(r.Code) == "SUCCESS"
}

func (r Response) isFailed() bool {
	return !r.isSuccess()
}
