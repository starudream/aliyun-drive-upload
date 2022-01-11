package aliyunDriveUpload

import (
	"net/http"
	"time"

	"github.com/go-sdk/lib/cache"
	"github.com/go-sdk/lib/conf"
	"github.com/go-sdk/lib/consts"
	"github.com/go-sdk/lib/httpx"
	"github.com/go-sdk/lib/rdx"
)

const (
	BaseURL = "https://auth.aliyundrive.com"

	AccessTokenPrefix  = "at:"
	RefreshTokenPrefix = "rt:"

	xEmptyContentType = "x-empty-content-type"
)

var (
	debug = conf.Get("debug").Bool()

	hc = httpx.New(httpx.WithDebug(debug)).SetPreRequestHook(hcHook)

	mc cache.Cache
)

func init() {
	if rdx.Default() == nil {
		mc = cache.NewMemoryCacheWithCleaner(time.Hour, time.Second, nil)
	} else {
		mc = cache.NewRedisCache(time.Hour, rdx.Default())
	}
}

func hcHook(_ *httpx.Client, req *http.Request) error {
	if req.Header.Get(xEmptyContentType) != "" {
		req.Header.Del(xEmptyContentType)
		req.Header.Set(consts.ContentType, "")
	}
	return nil
}
