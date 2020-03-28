package wechat

import (
	"context"
	"encoding/json"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/puper/wechat-sdk-server/errors"
	"github.com/puper/wechat-sdk-server/pkg/wechat/api"
)

type AccessToken struct {
	AppId       string `json:"appId,omitempty"`
	AccessToken string `json:"accessToken,omitempty"`
	Rev         int64  `json:"-"`
	State       int    `json:"-"`
	ExpiresAt   int64  `json:"expiresAt,omitempty"`
}

func (this *AccessToken) Clone() *AccessToken {
	if this != nil {
		reply := *this
		return &reply
	}
	return nil
}

func (this *AccessToken) IsInvalid() bool {
	return time.Now().Unix() >= this.ExpiresAt || this.State == StateInvalid
}

func (this *AccessToken) IsNeedRefresh() bool {
	return this.State == StateNeedRefresh
}

func (this *AccessToken) Decode(b []byte) error {
	return json.Unmarshal(b, this)
}

func (this *AccessToken) Encode() []byte {
	b, _ := json.Marshal(this)
	return b
}

type RefreshAccessToken struct {
	AppId           string
	LastAccessToken string
}

func (this *RefreshAccessToken) Decode(b []byte) error {
	return json.Unmarshal(b, this)
}

func (this *RefreshAccessToken) Encode() []byte {
	b, _ := json.Marshal(this)
	return b
}

type AccessTokenAccessAt struct {
	AppId    string
	AccessAt int64
}

func (this *AccessTokenAccessAt) Decode(b []byte) error {
	return json.Unmarshal(b, this)
}

func (this *AccessTokenAccessAt) Encode() []byte {
	b, _ := json.Marshal(this)
	return b
}

type GetAccessTokenArg struct {
	AppId              string
	CompareAccessToken string
	//MarkInvalid        bool
}

type GetAccessTokenReply = AccessToken

func (this *Wechat) GetAccessToken(arg *GetAccessTokenArg) (*GetAccessTokenReply, error) {
	_, err := this.GetApp(&GetAppArg{
		AppId: arg.AppId,
	})
	if err != nil {
		return nil, errors.ParamError(ErrAppNotFound).Log()
	}
	accessAt := int64(0)
	this.accessAtMutex.RLock()
	if _, ok := this.accessAts[arg.AppId]; ok {
		accessAt = this.accessAts[arg.AppId].AccessAt
	}
	this.accessAtMutex.RUnlock()
	if time.Now().Unix()-accessAt > this.config.AccessAtUpdateTime/1e3 {
		go func() {
			this.accessAtMutex.Lock()
			accessAt, ok := this.accessAts[arg.AppId]
			if !ok || time.Now().Unix()-accessAt.AccessAt > this.config.AccessAtUpdateTime/1e3 {
				accessAt := &AccessTokenAccessAt{
					AppId:    arg.AppId,
					AccessAt: time.Now().Unix(),
				}
				this.accessAts[accessAt.AppId] = accessAt
				this.accessAtMutex.Unlock()
				this.config.EtcdCli.Put(context.TODO(), GetKey(PrefixAccessAt, arg.AppId), string(accessAt.Encode()))
			} else {
				this.accessAtMutex.Unlock()
			}
		}()
	}
	this.accesstokensMutex.RLock()
	curToken, ok := this.accesstokens[arg.AppId]
	curToken = curToken.Clone()
	this.accesstokensMutex.RUnlock()
	if !ok {
		return this.getAccessToken(arg)
	}
	if curToken.AccessToken == arg.CompareAccessToken {
		return this.getAccessToken(arg)
	}
	if curToken.IsInvalid() {
		return this.getAccessToken(arg)
	}
	if curToken.IsNeedRefresh() {
		go this.getAccessToken(arg)
	}
	return curToken, nil
}

func (this *Wechat) getAccessToken(arg *GetAccessTokenArg) (*GetAccessTokenReply, error) {
	accesstoken, err := this.getAccessTokenFromEtcd(this.config.EtcdCli, arg.AppId)
	if err != nil {
		return nil, errors.ServerError(err).SetDetails(arg).LogWithTrace()
	}
	if accesstoken != nil &&
		!accesstoken.IsInvalid() &&
		accesstoken.AccessToken != arg.CompareAccessToken {
		return accesstoken.Clone(), nil
	}
	s, err := concurrency.NewSession(this.config.EtcdCli)
	if err != nil {
		return nil, errors.ServerError(err).SetDetails(arg).LogWithTrace()
	}
	defer s.Close()
	lockKey := GetKey(PrefixLock, PrefixAccessToken, arg.AppId)
	m := concurrency.NewMutex(s, lockKey)
	if err := m.Lock(context.TODO()); err != nil {
		return nil, errors.ServerError(err).SetDetails(arg).LogWithTrace()
	}
	defer m.Unlock(context.TODO())
	accesstoken, err = this.getAccessTokenFromEtcd(this.config.EtcdCli, arg.AppId)
	if err != nil {
		return nil, errors.ServerError(err).SetDetails(arg).LogWithTrace()
	}
	if accesstoken != nil &&
		!accesstoken.IsInvalid() &&
		accesstoken.AccessToken != arg.CompareAccessToken {
		return accesstoken.Clone(), nil
	}
	app, err := this.GetApp(&GetAppArg{
		AppId: arg.AppId,
	})
	if err != nil {
		return nil, err
	}
	reply, err := api.GetAccessToken(
		&api.GetAccessTokenArg{
			Client: this.config.RestyCli,
			Domain: this.config.ApiDomain,
			AppId:  app.AppId,
			Secret: app.Secret,
		},
	)
	if err != nil {
		return nil, errors.ServerError(err).SetDetails(arg).LogWithTrace()
	}

	accesstoken = &AccessToken{
		AppId:       arg.AppId,
		AccessToken: reply.AccessToken,
		ExpiresAt:   time.Now().Unix() + reply.ExpiresIn,
		Rev:         m.Header().GetRevision(),
	}
	refreshAccessToken := &RefreshAccessToken{
		AppId:           accesstoken.AppId,
		LastAccessToken: accesstoken.AccessToken,
	}
	ttl := int64(float64(reply.ExpiresIn) * 0.8)
	if ttl <= 0 {
		return accesstoken, nil
	}
	lease, err := this.config.EtcdCli.Grant(context.TODO(), ttl)
	if err != nil {
		accesstoken.State = StateInvalid
		this.accesstokensMutex.Lock()
		this.accesstokens[accesstoken.AppId] = accesstoken
		this.accesstokensMutex.Unlock()
		return accesstoken.Clone(), nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	_, err = this.config.EtcdCli.Txn(ctx).
		Then(
			clientv3.OpPut(GetKey(PrefixAccessToken, accesstoken.AppId), string(accesstoken.Encode())),
			clientv3.OpPut(
				GetKey(PrefixRefreshAccessToken, accesstoken.AppId),
				string(refreshAccessToken.Encode()),
				clientv3.WithLease(lease.ID),
			),
		).Commit()
	cancel()
	if err != nil {
		errors.ServerError(err).SetDetails(arg).LogWithTrace()
		accesstoken.State = StateInvalid
		this.accesstokensMutex.Lock()
		this.accesstokens[accesstoken.AppId] = accesstoken
		this.accesstokensMutex.Unlock()
		return accesstoken.Clone(), nil
	}
	this.accesstokensMutex.Lock()
	this.accesstokens[accesstoken.AppId] = accesstoken
	this.accesstokensMutex.Unlock()
	return accesstoken.Clone(), nil
}

func (this *Wechat) getAccessTokenFromEtcd(etcdKv clientv3.KV, appId string) (*AccessToken, error) {
	wr, err := etcdKv.Get(context.TODO(), GetKey(PrefixAccessToken, appId))
	if err != nil {
		return nil, errors.ServerError(err).SetDetails(appId).LogWithTrace()
	}
	for _, kv := range wr.Kvs {
		reply := &AccessToken{}
		if err := reply.Decode(kv.Value); err == nil {
			reply.Rev = wr.Header.GetRevision()
			return reply, nil
		}
		return nil, nil
	}
	return nil, nil
}
