package etcdcli

import (
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/puper/ppgo/v2/engine"
)

func Build(cfgKey string) func(*engine.Engine) (interface{}, error) {
	return func(e *engine.Engine) (interface{}, error) {
		cfg := &Config{}
		err := e.GetConfig().UnmarshalKey(cfgKey, cfg)
		if err != nil {
			return nil, err
		}
		cli, err := clientv3.New(
			clientv3.Config{
				Endpoints:   cfg.Endpoints,
				DialTimeout: time.Duration(cfg.DialTimeout) * time.Millisecond,
			},
		)
		return cli, err
	}
}

type Config struct {
	Endpoints   []string `json:"endpoints"`
	DialTimeout int64    `json:"dialTimeout"`
}
