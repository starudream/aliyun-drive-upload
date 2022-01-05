package aliyunDriveUpload

import (
	"github.com/go-sdk/lib/conf"
	"github.com/go-sdk/lib/httpx"
)

var RefreshToken = conf.Get("refresh_token").String()

func init() {
	hc = httpx.New(httpx.WithDebug(true)).SetPreRequestHook(hcHook)
}
