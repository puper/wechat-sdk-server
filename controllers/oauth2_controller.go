package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/puper/wechat-sdk-server/app"
	"github.com/puper/wechat-sdk-server/pkg/wechat"
	"github.com/puper/wechat-sdk-server/response"
)

type oauth2Controller struct{}

func (oauth2Controller) GetAuthUrl(ctx iris.Context) {
	resp := response.New(ctx)
	arg := &wechat.GetAuthUrlArg{}
	if err := ctx.ReadJSON(arg); err != nil {
		resp.ParamError().Send()
		return
	}
	reply, err := app.GetWechat().GetAuthUrl(arg)
	if err != nil {
		resp.Error(err).Send()
		return
	}
	resp.Result(reply).Send()

}

func (oauth2Controller) GetUserAccessToken(ctx iris.Context) {
	resp := response.New(ctx)
	arg := &wechat.GetUserAccessTokenArg{}
	if err := ctx.ReadJSON(arg); err != nil {
		resp.ParamError().Send()
		return
	}
	reply, err := app.GetWechat().GetUserAccessToken(arg)
	if err != nil {
		resp.Error(err).Send()
		return
	}
	resp.Result(reply).Send()

}

func (oauth2Controller) RefreshUserAccessToken(ctx iris.Context) {
	resp := response.New(ctx)
	arg := &wechat.RefreshUserAccessTokenArg{}
	if err := ctx.ReadJSON(arg); err != nil {
		resp.ParamError().Send()
		return
	}
	reply, err := app.GetWechat().RefreshUserAccessToken(arg)
	if err != nil {
		resp.Error(err).Send()
		return
	}
	resp.Result(reply).Send()

}

func (oauth2Controller) GetUserInfo(ctx iris.Context) {
	resp := response.New(ctx)
	arg := &wechat.GetUserInfoArg{}
	if err := ctx.ReadJSON(arg); err != nil {
		resp.ParamError().Send()
		return
	}
	reply, err := app.GetWechat().GetUserInfo(arg)
	if err != nil {
		resp.Error(err).Send()
		return
	}
	resp.Result(reply).Send()

}

func (oauth2Controller) CheckUserAccessToken(ctx iris.Context) {
	resp := response.New(ctx)
	arg := &wechat.CheckUserAccessTokenArg{}
	if err := ctx.ReadJSON(arg); err != nil {
		resp.ParamError().Send()
		return
	}
	reply, err := app.GetWechat().CheckUserAccessToken(arg)
	if err != nil {
		resp.Error(err).Send()
		return
	}
	resp.Result(reply).Send()

}
