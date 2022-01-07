package aliyunDriveUpload

import (
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/go-sdk/lib/codec/xml"
	"github.com/go-sdk/lib/consts"
	"github.com/go-sdk/lib/log"
)

func UploadFile(refreshToken, directory, filename string) (*CompleteResp, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	fs, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if fs.IsDir() {
		return nil, fmt.Errorf("upload: only support a file not a directory")
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
		Name:          fs.Name(),
		Type:          "file",
		CheckNameMode: "auto_rename",
		Size:          fs.Size(),
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

	// Method 1: use resty, but need comment `github.com/go-resty/resty/v2@v2.7.0/middleware.go:209` `getBodyCopy` function.

	// uploadResp, err := hc.
	// 	NewRequest().
	// 	SetHeader(emptyContentType, "true").
	// 	SetBody(&ProgressReader{reader: file, hook: pHook, totalBytes: fs.Size()}).
	// 	SetError(OSSResp{}).
	// 	Put(fileRes.PartInfoList[0].UploadUrl)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// if uploadResp.IsError() {
	// 	if e, ok := uploadResp.Error().(*OSSResp); ok {
	// 		return nil, fmt.Errorf("upload: %s, %s", e.Code, e.Message)
	// 	}
	// }

	// Method 2: use stdlib http client

	err = upload(file, fs, fileRes.PartInfoList[0].UploadUrl)
	if err != nil {
		return nil, err
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

func upload(file *os.File, fs os.FileInfo, url string) error {
	req, err := http.NewRequest(http.MethodPut, url, &ProgressReader{reader: file, hook: pHook, totalBytes: fs.Size()})
	if err != nil {
		return err
	}
	resp, err := hc.GetClient().Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode > 399 {
		ossResp := &OSSResp{}
		err = xml.NewDecoder(resp.Body).Decode(ossResp)
		if err != nil {
			return err
		}
		return fmt.Errorf("upload: %s, %s", ossResp.Code, ossResp.Message)
	}
	return nil
}

func pHook(event *ProgressEvent) {
	switch event.EventType {
	case transferStartedEvent:
		log.Debug("upload start")
	case transferDataEvent:
		if pSend() {
			log.Debugf("uploading, %.02f%%", float64(event.ConsumedBytes*100)/float64(event.TotalBytes))
		}
	case transferCompletedEvent:
		log.Info("upload success")
	case transferFailedEvent:
		log.Error("upload fail")
	}
}

var pUnix = int64(0)

func pSend() bool {
	x := time.Now().UnixMilli() / 500
	return x != atomic.SwapInt64(&pUnix, x)
}
