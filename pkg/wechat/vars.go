package wechat

const (
	Prefix                   = "wechat"
	PrefixApp                = "wechat/app"
	PrefixAccessToken        = "wechat/accesstoken"
	PrefixRefreshAccessToken = "wechat/refreshaccesstoken"
	PrefixAccessAt           = "wechat/accessat"

	PrefixLock = "lock"

	StateNormal          = 0 // 正常
	StateInvalid         = 1 // 上报失效. 不能使用
	StateNeedRefresh     = 2
	StateUpdatedDirectly = 3 // 手动更新的，更新检测版本
)
