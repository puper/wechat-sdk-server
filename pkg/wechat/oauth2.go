package wechat

import (
	"github.com/puper/wechat-sdk-server/pkg/wechat/api"
)

type GetAuthUrlArg = api.GetAuthUrlArg
type GetAuthUrlReply = api.GetAuthUrlReply

func (this *Wechat) GetAuthUrl(arg *GetAuthUrlArg) (*GetAuthUrlReply, error) {
	return api.GetAuthUrl(arg)
}

type GetUserAccessTokenArg = api.GetUserAccessTokenArg

type GetUserAccessTokenReply = api.GetUserAccessTokenReply

func (this *Wechat) GetUserAccessToken(arg *GetUserAccessTokenArg) (reply *GetUserAccessTokenReply, err error) {
	arg.Client = this.config.RestyCli
	return api.GetUserAccessToken(arg)
}

type RefreshUserAccessTokenArg = api.RefreshUserAccessTokenArg

type RefreshUserAccessTokenReply = api.RefreshUserAccessTokenReply

func (this *Wechat) RefreshUserAccessToken(arg *RefreshUserAccessTokenArg) (reply *RefreshUserAccessTokenReply, err error) {
	arg.Client = this.config.RestyCli
	return api.RefreshUserAccessToken(arg)
}

type GetUserInfoArg = api.GetUserInfoArg

type GetUserInfoReply = api.GetUserInfoReply

func (this *Wechat) GetUserInfo(arg *GetUserInfoArg) (reply *GetUserInfoReply, err error) {
	arg.Client = this.config.RestyCli
	return api.GetUserInfo(arg)
}

type CheckUserAccessTokenArg = api.CheckUserAccessTokenArg

type CheckUserAccessTokenReply = api.CheckUserAccessTokenReply

func (this *Wechat) CheckUserAccessToken(arg *CheckUserAccessTokenArg) (reply *CheckUserAccessTokenReply, err error) {
	arg.Client = this.config.RestyCli
	return api.CheckUserAccessToken(arg)
}
