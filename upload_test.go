package aliyunDriveUpload

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadFile(t *testing.T) {
	if RefreshToken == "" {
		t.SkipNow()
	}
	uploadResp, err := UploadFile(RefreshToken, "", "README.md")
	assert.NoError(t, err)

	downloadResp, err := GetDownloadURL(RefreshToken, uploadResp.FileId)
	assert.NoError(t, err)

	t.Log(downloadResp.Url)
	t.Log(downloadResp.RateLimit)
}
