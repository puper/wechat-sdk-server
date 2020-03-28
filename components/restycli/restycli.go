package restycli

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/puper/ppgo/v2/engine"
)

func Build(cfgKey string) func(*engine.Engine) (interface{}, error) {
	return func(e *engine.Engine) (interface{}, error) {
		cfg := &Config{}
		err := e.GetConfig().UnmarshalKey(cfgKey, cfg)
		if err != nil {
			return nil, err
		}
		hc := &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Millisecond,
		}
		transport := &http.Transport{
			MaxIdleConns:    cfg.MaxIdleConns,
			IdleConnTimeout: time.Duration(cfg.IdleConnTimeout) * time.Millisecond,
		}
		if cfg.InsecureSkipVerify {
			transport.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		hc.Transport = transport
		cli := resty.NewWithClient(hc)
		return cli, nil
	}
}

type Config struct {
	Timeout            int64 `json:"timeout"`
	MaxIdleConns       int   `json:"maxIdleConns"`
	IdleConnTimeout    int64 `json:"idleConnTimeout"`
	InsecureSkipVerify bool  `json:"insecureSkipVerify"`
}
