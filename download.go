package aliyunDriveUpload

import (
	"fmt"

	"github.com/go-sdk/lib/consts"
)

func GetDownloadURL(refreshToken, fileId string) (*DownloadResp, error) {
	token, err := GetToken(refreshToken)
	if err != nil {
		return nil, err
	}

	downloadResp, err := hc.
		NewRequest().
		SetHeader(consts.ContentType, consts.ContentTypeJSON).
		SetHeader(consts.Authorization, token.TokenType+" "+token.AccessToken).
		SetBody(DownloadReq{DriveId: token.DefaultDriveId, FileId: fileId, ExpireSec: 3600}).
		SetResult(DownloadResp{}).
		SetError(ErrResp{}).
		Post(BaseURL + "/v2/file/get_download_url")
	if err != nil {
		return nil, err
	}

	if e, ok := downloadResp.Error().(*ErrResp); ok {
		return nil, fmt.Errorf("download: %s, %s", e.Code, e.Message)
	}

	return downloadResp.Result().(*DownloadResp), nil
}
