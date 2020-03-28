package wechat

import (
	"context"
	"encoding/json"
	"time"

	"github.com/coreos/etcd/clientv3"
)

type App struct {
	AppId  string `json:"appId,omitempty"`
	Secret string `json:"secret,omitempty"`
}

func (this *App) Decode(b []byte) error {
	return json.Unmarshal(b, this)
}

func (this *App) Encode() []byte {
	b, _ := json.Marshal(this)
	return b
}

type CreateAppArg = App
type CreateAppReply struct{}

func (this *Wechat) CreateApp(arg *CreateAppArg) (*CreateAppReply, error) {
	_, err := this.config.EtcdCli.Put(context.TODO(), GetKey(PrefixApp, arg.AppId), string(arg.Encode()))
	if err != nil {
		return nil, err
	}
	return &CreateAppReply{}, nil
}

type UpdateAppArg = App
type UpdateAppReply struct{}

func (this *Wechat) UpdateApp(arg *UpdateAppArg) (*UpdateAppReply, error) {
	_, err := this.config.EtcdCli.Put(context.TODO(), GetKey(PrefixApp, arg.AppId), string(arg.Encode()))
	if err != nil {
		return nil, err
	}
	return &UpdateAppReply{}, nil
}

type DeleteAppArg struct {
	AppId string
}

type DeleteAppReply struct{}

func (this *Wechat) DeleteApp(arg *DeleteAppArg) (*DeleteAppReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	_, err := this.config.EtcdCli.Txn(ctx).
		Then(
			clientv3.OpDelete(GetKey(PrefixApp, arg.AppId)),
			clientv3.OpDelete(GetKey(PrefixAccessAt, arg.AppId)),
			clientv3.OpDelete(GetKey(PrefixAccessToken, arg.AppId)),
		).Commit()
	cancel()
	if err != nil {
		return nil, err
	}
	return &DeleteAppReply{}, nil
}

type GetAppArg struct {
	AppId string
}

type GetAppReply = App

func (this *Wechat) GetApp(arg *GetAppArg) (*GetAppReply, error) {
	this.appsMutex.RLock()
	defer this.appsMutex.RUnlock()
	app, ok := this.apps[arg.AppId]
	if ok {
		return app, nil
	}
	return nil, ErrAppNotFound
}
