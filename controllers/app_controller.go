package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/puper/wechat-sdk-server/app"
	"github.com/puper/wechat-sdk-server/pkg/wechat"
	"github.com/puper/wechat-sdk-server/response"
)

type appController struct{}

func (appController) Create(ctx iris.Context) {
	resp := response.New(ctx)
	arg := &wechat.CreateAppArg{}
	if err := ctx.ReadJSON(arg); err != nil {
		resp.ParamError().Send()
		return
	}
	reply, err := app.GetWechat().CreateApp(arg)
	if err != nil {
		resp.Error(err).Send()
		return
	}
	resp.Result(reply).Send()

}

func (appController) Update(ctx iris.Context) {
	resp := response.New(ctx)
	arg := &wechat.UpdateAppArg{}
	if err := ctx.ReadJSON(arg); err != nil {
		resp.ParamError().Send()
		return
	}
	reply, err := app.GetWechat().UpdateApp(arg)
	if err != nil {
		resp.Error(err).Send()
		return
	}
	resp.Result(reply).Send()
}

func (appController) Delete(ctx iris.Context) {
	resp := response.New(ctx)
	arg := &wechat.DeleteAppArg{}
	if err := ctx.ReadJSON(arg); err != nil {
		resp.ParamError().Send()
		return
	}
	reply, err := app.GetWechat().DeleteApp(arg)
	if err != nil {
		resp.Error(err).Send()
		return
	}
	resp.Result(reply).Send()
}

func (appController) Get(ctx iris.Context) {
	resp := response.New(ctx)
	arg := &wechat.GetAppArg{}
	if err := ctx.ReadJSON(arg); err != nil {
		resp.ParamError().Send()
		return
	}
	reply, err := app.GetWechat().GetApp(arg)
	if err != nil {
		resp.Error(err).Send()
		return
	}
	resp.Result(reply).Send()
}
