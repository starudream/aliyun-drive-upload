package aliyunDriveUpload

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetToken(t *testing.T) {
	if RefreshToken == "" {
		t.SkipNow()
	}
	_, err := GetToken(RefreshToken)
	assert.NoError(t, err)
}
