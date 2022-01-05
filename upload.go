package aliyunDriveUpload

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-sdk/lib/consts"
)

func UploadFile(refreshToken, directory, filename string) (*CompleteResp, error) {
	bs, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	token, err := GetToken(refreshToken)
	if err != nil {
		return nil, err
	}

	if directory == "" {
		directory = "root"
	}

	fileReq := FileReq{
		DriveId:       token.DefaultDriveId,
		PartInfoList:  []FilePartReq{{PartNumber: 1}},
		ParentFileId:  directory,
		Name:          filepath.Base(filename),
		Type:          "file",
		CheckNameMode: "auto_rename",
		Size:          len(bs),
	}

	fileResp, err := hc.
		NewRequest().
		SetHeader(consts.ContentType, consts.ContentTypeJSON).
		SetHeader(consts.Authorization, token.TokenType+" "+token.AccessToken).
		SetBody(fileReq).
		SetResult(FileResp{}).
		SetError(ErrResp{}).
		Post(BaseURL + "/v2/file/create_with_proof")
	if err != nil {
		return nil, err
	}

	if e, ok := fileResp.Error().(*ErrResp); ok {
		return nil, fmt.Errorf("upload: %s, %s", e.Code, e.Message)
	}

	fileRes := fileResp.Result().(*FileResp)

	if len(fileRes.PartInfoList) != 1 {
		return nil, fmt.Errorf("upload: only support one part file")
	}

	uploadResp, err := hc.
		NewRequest().
		SetHeader(emptyContentType, "true").
		SetBody(bs).
		SetError(OSSResp{}).
		Put(fileRes.PartInfoList[0].UploadUrl)
	if err != nil {
		return nil, err
	}

	if uploadResp.IsError() {
		if e, ok := uploadResp.Error().(*OSSResp); ok {
			return nil, fmt.Errorf("upload: %s, %s", e.Code, e.Message)
		}
	}

	completeResp, err := hc.
		NewRequest().
		SetHeader(consts.ContentType, consts.ContentTypeJSON).
		SetHeader(consts.Authorization, token.TokenType+" "+token.AccessToken).
		SetBody(CompleteReq{DriveId: fileRes.DriveId, UploadId: fileRes.UploadId, FileId: fileRes.FileId}).
		SetResult(CompleteResp{}).
		SetError(ErrResp{}).
		Post(BaseURL + "/v2/file/complete")
	if err != nil {
		return nil, err
	}

	if e, ok := completeResp.Error().(*ErrResp); ok {
		return nil, fmt.Errorf("upload: %s, %s", e.Code, e.Message)
	}

	return completeResp.Result().(*CompleteResp), nil
}