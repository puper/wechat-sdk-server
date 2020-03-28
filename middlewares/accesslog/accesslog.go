package accesslog

import (
	"fmt"
	"strconv"
	"time"

	"github.com/tidwall/gjson"

	"github.com/sirupsen/logrus"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

func New(logger *logrus.Logger) context.Handler {
	return func(ctx iris.Context) {
		var status, ip, method, path string
		var latency time.Duration
		var startTime, endTime time.Time
		startTime = time.Now()

		ctx.Record()
		ctx.Next()

		endTime = time.Now()
		latency = endTime.Sub(startTime)
		status = strconv.Itoa(ctx.GetStatusCode())
		ip = ctx.RemoteAddr()
		method = ctx.Method()
		path = ctx.Request().URL.RequestURI()
		resp := gjson.ParseBytes(ctx.Recorder().Body())

		logger.WithFields(logrus.Fields{
			"start":   startTime,
			"status":  status,
			"latency": fmt.Sprintf("%4v", latency),
			"ip":      ip,
			"method":  method,
			"path":    path,
			"code":    resp.Get("code").String(),
		}).Infoln(resp.Get("message").String())
	}
}
