package routes

import (
	"github.com/puper/wechat-sdk-server/app"
	"github.com/puper/wechat-sdk-server/controllers"

	"github.com/puper/wechat-sdk-server/middlewares/accesslog"

	"github.com/kataras/iris/v12"
)

func Configure(router *iris.Application) {
	router.Use(accesslog.New(app.GetLog("access")))
	v1 := router.Party("/api/v1")
	{
		v1.Post("/apps/create", controllers.AppController.Create)
		v1.Post("/apps/update", controllers.AppController.Update)
		v1.Post("/apps/delete", controllers.AppController.Delete)
		v1.Post("/apps/get", controllers.AppController.Get)

		v1.Post("/accesstokens/get", controllers.AccessTokenController.Get)

		v1.Post("/oauth2/authurls/get", controllers.Oauth2Cotroller.GetAuthUrl)
		v1.Post("/oauth2/userinfos/get", controllers.Oauth2Cotroller.GetUserInfo)
		v1.Post("/oauth2/accesstokens/get", controllers.Oauth2Cotroller.GetUserAccessToken)
		v1.Post("/oauth2/accesstoekns/refresh", controllers.Oauth2Cotroller.RefreshUserAccessToken)
		v1.Post("/oauth2/accesstokens/check", controllers.Oauth2Cotroller.CheckUserAccessToken)

	}
}
