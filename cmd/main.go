package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-sdk/lib/app"
	"github.com/go-sdk/lib/flag"
	"github.com/go-sdk/lib/log"
	"github.com/go-sdk/lib/timex"

	aliyunDriveUpload "github.com/starudream/aliyun-drive-upload"
)

const link = "https://passport.aliyundrive.com/mini_login.htm?lang=zh_cn&appName=aliyun_drive&appEntrance=web&styleType=auto"

var (
	Ver   bool
	Token string
	Dir   string
	Files []string
	Fids  []string
)

func init() {
	flag.BoolVar(&Ver, "version", false, "show versions")
	flag.StringVar(&Token, "token", "", "refresh token\nopen this link and click F12 to open devtools, then login to your account.\nfind the request in the network and get the 'bizExt' value in response, copy 'refreshToken' in the base64 decrypted json.\nlink: "+link)
	flag.StringVar(&Dir, "dir", "root", "directory id of you upload\nopen the 'aliyundrive.com' and click into the directory, you can find id in the url.")
	flag.StringSliceVar(&Files, "file", []string{}, "local file path\nif multiple, use '-file a -file b'.")
	flag.StringSliceVar(&Fids, "fid", []string{}, "remote file id\nif multiple, use '-fid xxx -fid yyy'.")
	flag.Parse()

	if Ver {
		fmt.Println(app.VersionInfo())
		os.Exit(0)
	}

	Files = sliceTrimSpace(Files)
	Fids = sliceTrimSpace(Fids)

	if Token == "" || (len(Files) == 0 && len(Fids) == 0) {
		fmt.Println("missing args")
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	a := app.New("aliyun-drive-upload")
	defer a.Recover()

	a.Add(func() error {
		for i := 0; i < len(Files); i++ {
			l := log.WithField("file", Files[i])
			uploadResp, err := aliyunDriveUpload.UploadFile(Token, "", Files[i])
			if err != nil {
				l.Error(err)
				continue
			}
			l.Infof("upload success, fid: %s", uploadResp.FileId)
			downloadResp, err := aliyunDriveUpload.GetDownloadURL(Token, uploadResp.FileId)
			if err != nil {
				l.Error(err)
				continue
			}
			l.Infof("get download url success, url: %s, expiration: %s", downloadResp.Url, downloadResp.Expiration.Local().Format(timex.DateTime))
			if downloadResp.RateLimit.PartSpeed != -1 {
				l.Warnf("download speed is limited, %d", downloadResp.RateLimit.PartSpeed)
			}
		}
		return nil
	})

	a.Add(func() error {
		for i := 0; i < len(Fids); i++ {
			l := log.WithField("fid", Fids[i])
			downloadResp, err := aliyunDriveUpload.GetDownloadURL(Token, Fids[i])
			if err != nil {
				l.Error(err)
				continue
			}
			l.Infof("get download url success, url: %s, expiration: %s", downloadResp.Url, downloadResp.Expiration.Local().Format(timex.DateTime))
			if downloadResp.RateLimit.PartSpeed != -1 {
				l.Warnf("download speed is limited, %d", downloadResp.RateLimit.PartSpeed)
			}
		}
		return nil
	})

	err := a.Once()
	if err != nil {
		log.Fatal(err)
	}
}

func sliceTrimSpace(s1 []string) []string {
	var s2 []string
	for i := 0; i < len(s1); i++ {
		s := strings.TrimSpace(s1[i])
		if s == "" {
			continue
		}
		s2 = append(s2, s)
	}
	return s2
}
