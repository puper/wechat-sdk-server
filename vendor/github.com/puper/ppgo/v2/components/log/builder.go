package log

import (
	"github.com/puper/ppgo/helpers"
	"github.com/puper/ppgo/v2/engine"
)

func Builder(configKey string) engine.Builder {
	return func(e *engine.Engine) (interface{}, error) {
		cfg := e.GetConfig().Get(configKey)
		c := &Config{}
		if err := helpers.StructDecode(cfg, c, "json"); err != nil {
			return nil, err
		}
		return New(c)
	}
}
