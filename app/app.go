package app

import (
	"github.com/puper/ppgo/v2/components/log"
	"github.com/puper/ppgo/v2/engine"
	"github.com/puper/wechat-sdk-server/pkg/wechat"
	"github.com/sirupsen/logrus"
)

var (
	app *engine.Engine
)

func Set(e *engine.Engine) {
	app = e
}

func Get() *engine.Engine {
	return app
}

func GetConfig() *engine.Config {
	return app.GetConfig()
}

func GetLog(name string) *logrus.Logger {
	return app.Get("log").(*log.Log).Get(name)
}

func GetWechat() *wechat.Wechat {
	return app.Get("wechat").(*wechat.Wechat)
}
