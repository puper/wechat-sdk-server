package bootstrap

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/puper/wechat-sdk-server/app"
	"github.com/puper/wechat-sdk-server/components/etcdcli"
	"github.com/puper/wechat-sdk-server/components/restycli"
	"github.com/puper/wechat-sdk-server/components/wechat"
	"github.com/puper/wechat-sdk-server/routes"

	"github.com/sirupsen/logrus"

	"github.com/kataras/iris/v12"
	"github.com/puper/ppgo/v2/components/irisapp"
	"github.com/puper/ppgo/v2/components/log"
	"github.com/puper/ppgo/v2/engine"
	"github.com/spf13/viper"

	terrors "github.com/puper/ppgo/errors"
)

func Bootstrap(cfgFile string) error {
	conf := viper.New()
	conf.SetConfigFile(cfgFile)
	if err := conf.ReadInConfig(); err != nil {
		return err
	}
	e := engine.New(conf)
	app.Set(e)
	e.Register("log", log.Builder("log"))
	e.Register("etcdcli", etcdcli.Build("etcdcli"))
	e.Register("restycli", restycli.Build("restycli"))

	e.Register("wechat", wechat.Build("wechat"), "log", "etcdcli")
	e.Register("web", func(e *engine.Engine) (interface{}, error) {
		s := &http.Server{
			ReadTimeout:  time.Minute * 10,
			WriteTimeout: time.Minute * 5,
			IdleTimeout:  time.Minute,
			Addr:         e.GetConfig().GetString("web.addr"),
		}
		app := &irisapp.Application{
			Application: iris.New(),
		}
		// register routes here
		routes.Configure(app.Application)
		go app.Run(
			iris.Server(s),
			iris.WithoutServerError(
				iris.ErrServerClosed,
			),
			iris.WithoutPathCorrection,
		)
		return app, nil
	}, "log")
	terrors.DefaultLoggerFunc = func() *logrus.Logger {
		return app.GetLog("")
	}
	err := e.Build()
	if err != nil {
		return err
	}
	defer e.Close()
	stop := make(chan struct{})
	go func() {
		sChan := make(chan os.Signal)
		for {
			signal.Notify(sChan, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
			sig := <-sChan
			switch sig {
			case os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				stop <- struct{}{}
			}

		}
	}()
	<-stop
	return nil
}
