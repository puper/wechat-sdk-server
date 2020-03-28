package log

import (
	"github.com/sirupsen/logrus"
)

type Log struct {
	logs map[string]*logrus.Logger
}

type Config map[string]map[string]string

func New(cfg *Config) (*Log, error) {
	instance := &Log{
		logs: make(map[string]*logrus.Logger),
	}
	instance.logs["default"] = logrus.StandardLogger()
	for name, config := range *cfg {
		l := logrus.New()
		level, err := logrus.ParseLevel(config["level"])
		if err != nil {
			return nil, err
		}
		l.Level = level
		if config["format"] == "json" {
			l.Formatter = &logrus.JSONFormatter{}
		} else {
			l.Formatter = &logrus.TextFormatter{}
		}
		if config["out"] != "std" {
			arf, err := newAutoRotateFile(config["out"], RotateTypeDay)
			if err != nil {
				return nil, err
			}
			l.Out = arf
		}
		instance.logs[name] = l
	}
	return instance, nil
}

func (this *Log) Get(name string) *logrus.Logger {
	l, ok := this.logs[name]
	if ok {
		return l
	}
	return this.logs["default"]
}
