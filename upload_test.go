package aliyunDriveUpload

import (
	"testing"

	"github.com/go-sdk/lib/testx"
)

func TestUploadFile(t *testing.T) {
	if RefreshToken == "" {
		t.SkipNow()
	}
	uploadResp, err := UploadFile(RefreshToken, "", "README.md")
	testx.RequireNoError(t, err)

	downloadResp, err := GetDownloadURL(RefreshToken, uploadResp.FileId)
	testx.RequireNoError(t, err)

	t.Log(downloadResp.Url)
	t.Log(downloadResp.RateLimit)
}
