package aliyunDriveUpload

import (
	"testing"

	"github.com/go-sdk/lib/testx"
)

func TestGetToken(t *testing.T) {
	if RefreshToken == "" {
		t.SkipNow()
	}
	_, err := GetToken(RefreshToken)
	testx.AssertNoError(t, err)
}
