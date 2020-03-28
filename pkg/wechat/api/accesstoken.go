package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type GetAccessTokenArg struct {
	Client *resty.Client
	Domain string
	AppId  string
	Secret string
}

type GetAccessTokenReply struct {
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int64  `json:"expires_in,omitempty"`
}

func GetAccessToken(arg *GetAccessTokenArg) (reply *GetAccessTokenReply, err error) {
	response, err := arg.Client.R().
		Get(
			fmt.Sprintf(
				"https://%v/cgi-bin/token?grant_type=%v&appid=%v&secret=%v",
				arg.Domain,
				"client_credential",
				arg.AppId,
				arg.Secret,
			),
		)

	reply = &GetAccessTokenReply{}
	err = parseResponse(response, err, reply)
	return reply, err
}
