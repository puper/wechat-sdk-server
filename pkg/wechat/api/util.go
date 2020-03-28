package api

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

func parseResponse(response *resty.Response, err error, reply interface{}) error {
	if err != nil {
		return &ApiError{
			Type:    ErrTypeConn,
			Message: err.Error(),
		}
	}
	j := gjson.ParseBytes(response.Body())
	errCode := j.Get("errcode").Int()
	if errCode != Success {

		return &ApiError{
			Code:    j.Get("errcode").Int(),
			Message: j.Get("errmsg").String(),
		}
	}
	err = json.Unmarshal(response.Body(), reply)
	if err != nil {
		return &ApiError{
			Type:    ErrTypeResp,
			Message: err.Error(),
		}
	}
	return nil
}
