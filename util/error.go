package util

import "fmt"

type ErrorMsg struct {
	HttpRespCode int    `json:"http_resp_code,omitempty"`
	GrpcRespCode int    `json:"grpc_resp_code,omitempty"`
	Description  string `json:"description"`
}

func (em *ErrorMsg) Error() string {
	return fmt.Sprintf("GRPC ErrCode = %d HttpRespCode = %d Desc = %s", em.GrpcRespCode, em.HttpRespCode, em.Description)
}
