package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/puper/wechat-sdk-server/app"
	"github.com/puper/wechat-sdk-server/pkg/wechat"
	"github.com/puper/wechat-sdk-server/response"
)

type accessTokenController struct{}

func (accessTokenController) Get(ctx iris.Context) {
	resp := response.New(ctx)
	arg := new(wechat.GetAccessTokenArg)
	if err := ctx.ReadJSON(arg); err != nil {
		resp.ParamError().Send()
		return
	}
	reply, err := app.GetWechat().GetAccessToken(arg)
	if err != nil {
		resp.Error(err).Send()
		return
	}
	resp.Result(reply).Send()
}
