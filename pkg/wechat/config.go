package wechat

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type Config struct {
	ApiDomain                 string           `json:"apiDomain,omitempty"`
	AccessAtUpdateTime        int64            `json:"accessAtUpdateTime,omitempty"`
	DropAccessTokenNotUseTime int64            `json:"dropAccessTokenNotUseTime,omitempty"`
	EtcdCliName               string           `json:"etcdCliName,omitempty"`
	RestyCliName              string           `json:"restyCliName,omitempty"`
	Logger                    *logrus.Logger   `json:"-"`
	EtcdCli                   *clientv3.Client `json:"-"`
	RestyCli                  *resty.Client    `json:"-"`
}
