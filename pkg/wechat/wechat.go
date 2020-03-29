package wechat

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/puper/wechat-sdk-server/errors"
	"golang.org/x/net/context"
)

func New(cfg *Config) (*Wechat, error) {
	wc := &Wechat{
		apps:         map[string]*App{},
		accesstokens: map[string]*AccessToken{},
		accessAts:    map[string]*AccessTokenAccessAt{},
		cancelc:      make(chan struct{}),
		donec:        make(chan struct{}),

		config: cfg,
	}
	err := wc.start()
	if err != nil {
		return nil, err
	}
	return wc, nil
}

type Wechat struct {
	appsMutex         sync.RWMutex
	apps              map[string]*App
	accesstokensMutex sync.RWMutex
	accesstokens      map[string]*AccessToken
	accessAtMutex     sync.RWMutex
	accessAts         map[string]*AccessTokenAccessAt
	config            *Config

	cancelc chan struct{}
	donec   chan struct{}
}

func (this *Wechat) start() error {
	wr, err := this.config.EtcdCli.Get(context.TODO(), Prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	hasRefreshTokens := map[string]bool{}
	rev := wr.Header.GetRevision()
	for _, kv := range wr.Kvs {
		k := string(kv.Key)
		if strings.HasPrefix(k, GetKey(PrefixApp, "")) {
			app := &App{}
			if err := app.Decode(kv.Value); err == nil {
				this.apps[app.AppId] = app
			} else {
				// log error
			}
		} else if strings.HasPrefix(k, GetKey(PrefixAccessToken, "")) {
			accesstoken := &AccessToken{}
			if err := accesstoken.Decode(kv.Value); err == nil {
				accesstoken.Rev = rev
				this.accesstokens[accesstoken.AppId] = accesstoken
			} else {

			}
		} else if strings.HasPrefix(k, GetKey(PrefixAccessAt, "")) {
			accessAt := &AccessTokenAccessAt{}
			if err := accessAt.Decode(kv.Value); err == nil {
				this.accessAts[accessAt.AppId] = accessAt
			} else {

			}
		} else if strings.HasPrefix(k, GetKey(PrefixRefreshAccessToken, "")) {
			// 不存在定时器的，需要mark need update,说明已经过了自动刷新期
			refreshAccessToken := &RefreshAccessToken{}
			if err := refreshAccessToken.Decode(kv.Value); err == nil {
				hasRefreshTokens[refreshAccessToken.AppId] = true
			}
		} else {

		}
	}
	for _, v := range this.accesstokens {
		if !hasRefreshTokens[v.AppId] {
			v.State = StateNeedRefresh
		}
	}
	this.config.Logger.Infoln("load wechat data success")
	this.config.Logger.Infoln("watch wechat data start")
	go this.watchRetry(rev + 1)
	return nil
}

func (this *Wechat) watchRetry(rev int64) {
	defer close(this.donec)
	for {
		rev = this.watch(rev)
		select {
		case <-this.cancelc:
			return
		default:
			time.Sleep(time.Second * 5)
		}
	}
}

func (this *Wechat) watch(rev int64) int64 {
	ctx := context.TODO()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ctx = clientv3.WithRequireLeader(ctx)
	rch := this.config.EtcdCli.Watch(ctx, GetKey(Prefix, ""), clientv3.WithPrefix(), clientv3.WithRev(rev), clientv3.WithPrevKV())
	for {
		select {
		case wr, ok := <-rch:
			if !ok {
				return rev + 1
			}
			if wr.Err() != nil {
				return rev + 1
			}
			fmt.Println("watch rev", rev)
			if wr.Header.GetRevision() > 0 {
				rev = wr.Header.GetRevision()
			}
			for _, event := range wr.Events {
				k := string(event.Kv.Key)
				if event.Type == clientv3.EventTypePut {
					if strings.HasPrefix(k, GetKey(PrefixApp, "")) {
						app := &App{}
						if err := app.Decode(event.Kv.Value); err == nil {
							this.appsMutex.Lock()
							this.apps[app.AppId] = app
							this.appsMutex.Unlock()
						} else {
							errors.ServerError(err).SetDetails(event.Kv.Value).LogWithTrace()
						}
					} else if strings.HasPrefix(k, GetKey(PrefixAccessToken, "")) {
						accesstoken := &AccessToken{}
						if err := accesstoken.Decode(event.Kv.Value); err == nil {
							accesstoken.Rev = rev
							this.accesstokensMutex.Lock()
							prevToken, ok := this.accesstokens[accesstoken.AppId]
							if !ok || prevToken.Rev <= rev {
								this.accesstokens[accesstoken.AppId] = accesstoken
							}
							this.accesstokensMutex.Unlock()
						} else {
							errors.ServerError(err).SetDetails(event.Kv.Value).LogWithTrace()
						}
					} else if strings.HasPrefix(k, GetKey(PrefixAccessAt, "")) {
						accessAt := &AccessTokenAccessAt{}
						if err := accessAt.Decode(event.Kv.Value); err == nil {
							this.accessAtMutex.Lock()
							this.accessAts[accessAt.AppId] = accessAt
							this.accessAtMutex.Unlock()
						} else {
							errors.ServerError(err).SetDetails(event.Kv.Value).LogWithTrace()
						}
					}
				} else if event.Type == clientv3.EventTypeDelete {
					if strings.HasPrefix(k, GetKey(PrefixApp, "")) {
						appId := strings.Split(k, GetKey(PrefixApp, ""))[1]
						this.appsMutex.Lock()
						delete(this.apps, appId)
						this.appsMutex.Unlock()
					} else if strings.HasPrefix(k, GetKey(PrefixAccessToken, "")) {
						appId := strings.Split(k, GetKey(PrefixAccessToken, ""))[1]
						this.accesstokensMutex.Lock()
						delete(this.accesstokens, appId)
						this.accesstokensMutex.Unlock()
					} else if strings.HasPrefix(k, GetKey(PrefixAccessAt, "")) {
						appId := strings.Split(k, GetKey(PrefixAccessAt, ""))[1]
						this.accessAtMutex.Lock()
						delete(this.accessAts, appId)
						this.accessAtMutex.Unlock()
					} else if strings.HasPrefix(k, GetKey(PrefixRefreshAccessToken, "")) {
						if event.PrevKv == nil {
							continue
						}
						refreshAccessToken := &RefreshAccessToken{}
						if err := refreshAccessToken.Decode(event.PrevKv.Value); err == nil {
							this.accessAtMutex.RLock()
							accessAt := int64(0)
							if _, ok := this.accessAts[refreshAccessToken.AppId]; ok {
								accessAt = this.accessAts[refreshAccessToken.AppId].AccessAt
							}
							this.accessAtMutex.RUnlock()
							if ok {
								if time.Now().Unix()-accessAt <= this.config.DropAccessTokenNotUseTime/1e3 {
									go func() {
										time.Sleep(time.Millisecond * time.Duration(rand.Intn(10000)))
										this.getAccessToken(&GetAccessTokenArg{
											AppId:              refreshAccessToken.AppId,
											CompareAccessToken: refreshAccessToken.LastAccessToken,
										})
									}()
								}
							}
						} else {
							errors.ServerError(err).SetDetails(string(event.Kv.Value)).LogWithTrace()
						}
					}
				}
			}
		case <-this.cancelc:
			return rev + 1
		}
	}
}

func (this *Wechat) Close() error {
	close(this.cancelc)
	<-this.donec
	return nil
}
