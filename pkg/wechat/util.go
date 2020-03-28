package wechat

import (
	"strings"
)

func GetKey(parts ...string) string {
	return strings.Join(parts, "/")
}
