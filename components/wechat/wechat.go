package wechat

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/go-resty/resty/v2"
	"github.com/puper/ppgo/v2/components/log"
	"github.com/puper/ppgo/v2/engine"
	"github.com/puper/wechat-sdk-server/pkg/wechat"
)

func Build(cfgKey string) func(*engine.Engine) (interface{}, error) {
	return func(e *engine.Engine) (interface{}, error) {
		cfg := &wechat.Config{}
		err := e.GetConfig().UnmarshalKey(cfgKey, cfg)
		if err != nil {
			return nil, err
		}
		cfg.EtcdCli = e.Get(cfg.EtcdCliName).(*clientv3.Client)
		cfg.RestyCli = e.Get(cfg.RestyCliName).(*resty.Client)
		cfg.Logger = e.Get("log").(*log.Log).Get("")
		return wechat.New(cfg)
	}
}
