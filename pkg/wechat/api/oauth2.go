package api

import (
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
)

type GetAuthUrlArg struct {
	AppId       string `json:"appId"`
	RedirectUri string `json:"redirectUri"`
	Scope       string `json:"scope"`
	State       string `json:"state"`
}

type GetAuthUrlReply struct {
	Url string `json:"url"`
}

func GetAuthUrl(arg *GetAuthUrlArg) (*GetAuthUrlReply, error) {
	reply := &GetAuthUrlReply{
		Url: fmt.Sprintf(
			"https://open.weixin.qq.com/connect/oauth2/authorize?appid=%v&redirect_uri=%v&response_type=code&scope=%v&state=%v#wechat_redirect",
			arg.AppId,
			url.QueryEscape(arg.RedirectUri),
			arg.Scope,
			url.QueryEscape(arg.State),
		),
	}
	return reply, nil
}

type GetUserAccessTokenArg struct {
	Client *resty.Client `json:"-"`
	Domain string        `json:"domain,omitempty"`
	AppId  string        `json:"appId,omitempty"`
	Secret string        `json:"secret,omitempty"`
	Code   string        `json:"code,omitempty"`
}

type GetUserAccessTokenReply struct {
	AccessToken  string `json:"accessToken"`
	ExpiresIn    int64  `json:"expiresIn"`
	RefreshToken string `json:"refreshToken"`
	OpenId       string `json:"openId"`
	Scope        string `json:"scope"`
}

func GetUserAccessToken(arg *GetUserAccessTokenArg) (reply *GetUserAccessTokenReply, err error) {
	response, err := arg.Client.R().
		Get(
			fmt.Sprintf(
				"https://%v/sns/oauth2/access_token?appid=%v&secret=%v&code=%v&grant_type=authorization_code",
				arg.Domain,
				arg.AppId,
				arg.Secret,
				arg.Code,
			),
		)
	err = parseResponse(response, err, reply)
	return reply, err
}

type RefreshUserAccessTokenArg struct {
	Client       *resty.Client `json:"-"`
	Domain       string        `json:"domain"`
	AppId        string        `json:"appId"`
	RefreshToken string        `json:"refreshToken"`
}

type RefreshUserAccessTokenReply struct {
	AccessToken  string `json:"accessToken"`
	ExpiresIn    int64  `json:"expiresIn"`
	RefreshToken string `json:"refreshToken"`
	OpenId       string `json:"openId"`
	Scope        string `json:"scope"`
}

func RefreshUserAccessToken(arg *RefreshUserAccessTokenArg) (reply *RefreshUserAccessTokenReply, err error) {
	response, err := arg.Client.R().
		Get(
			fmt.Sprintf(
				"https://%v/sns/oauth2/refresh_token?appid=%v&grant_type=refresh_token&refresh_token=%v",
				arg.Domain,
				arg.AppId,
				arg.RefreshToken,
			),
		)
	err = parseResponse(response, err, reply)
	return reply, err
}

type GetUserInfoArg struct {
	Client      *resty.Client `json:"-"`
	Domain      string        `json:"domain"`
	AccessToken string        `json:"accessToken"`
	OpenId      string        `json:"openId"`
}

type GetUserInfoReply struct {
	OpenId     string `json:"openId"`
	Nickname   string `json:"nickname"`
	HeadImgUrl string `json:"headImgUrl"`
	UnionId    string `json:"unionId"`
}

func GetUserInfo(arg *GetUserInfoArg) (reply *GetUserInfoReply, err error) {
	response, err := arg.Client.R().
		Get(
			fmt.Sprintf(
				"https://%v/sns/userinfo?access_token=%v&openid=%v&lang=zh_CN",
				arg.Domain,
				arg.AccessToken,
				arg.OpenId,
			),
		)
	err = parseResponse(response, err, reply)
	return reply, err
}

type CheckUserAccessTokenArg struct {
	Client      *resty.Client `json:"-"`
	Domain      string        `json:"domain"`
	AccessToken string        `json:"accessToken"`
	OpenId      string        `json:"openId"`
}

type CheckUserAccessTokenReply struct {
	Valid bool `json:"valid"`
}

func CheckUserAccessToken(arg *CheckUserAccessTokenArg) (reply *CheckUserAccessTokenReply, err error) {
	response, err := arg.Client.R().
		Get(
			fmt.Sprintf(
				"https://%v/sns/auth?access_token=%v&openid=%v",
				arg.Domain,
				arg.AccessToken,
				arg.OpenId,
			),
		)
	err = parseResponse(response, err, reply)
	if apiErr, ok := err.(*ApiError); ok && apiErr.Type == ErrTypeApi && apiErr.Code == InvalidOpenId {
		return reply, nil
	}
	return reply, err
}
