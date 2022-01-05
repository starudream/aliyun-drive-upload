package aliyunDriveUpload

import (
	"github.com/go-sdk/lib/conf"
)

var RefreshToken = conf.Get("refresh_token").String()

func init() {
	hc = hc.SetDebug(true)
}
