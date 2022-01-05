package aliyunDriveUpload

import (
	"net/http"
	"time"

	"github.com/go-sdk/lib/cache"
	"github.com/go-sdk/lib/conf"
	"github.com/go-sdk/lib/consts"
	"github.com/go-sdk/lib/httpx"
)

const (
	BaseURL = "https://auth.aliyundrive.com"

	AccessTokenPrefix  = "at:"
	RefreshTokenPrefix = "rt:"

	emptyContentType = "x-empty-content-type"
)

var (
	debug = conf.Get("debug").Bool()

	hc = httpx.New(httpx.WithDebug(debug)).SetPreRequestHook(hcHook)

	mc = cache.NewMemoryCacheWithCleaner(time.Hour, time.Second, nil)
)

func hcHook(_ *httpx.Client, req *http.Request) error {
	if req.Header.Get(emptyContentType) != "" {
		req.Header.Del(emptyContentType)
		req.Header.Set(consts.ContentType, "")
	}
	return nil
}
