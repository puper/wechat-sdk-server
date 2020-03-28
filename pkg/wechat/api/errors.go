package api

import "fmt"

const (
	Success     = 0
	ErrTypeConn = "connerr"
	ErrTypeResp = "resperr"
	ErrTypeApi  = "apierr"
)

type ApiError struct {
	Type    string
	Code    int64
	Message string
}

func (this *ApiError) Error() string {
	return fmt.Sprintf("%v:%v", this.Code, this.Message)
}
