package aliyunDriveUpload

import (
	"fmt"
	"time"

	"github.com/go-sdk/lib/consts"
)

func GetToken(refreshToken string) (*TokenResp, error) {
	at, exist, err := mc.Get(AccessTokenPrefix + refreshToken)
	if err != nil {
		return nil, err
	}
	if exist {
		return at.(*TokenResp), nil
	}

	rt, exist, err := mc.Get(RefreshTokenPrefix + refreshToken)
	if err != nil {
		return nil, err
	}
	if exist {
		refreshToken = rt.(string)
	}

	resp, err := hc.
		NewRequest().
		SetHeader(consts.ContentType, consts.ContentTypeJSON).
		SetBody(TokenReq{GrantType: "refresh_token", RefreshToken: refreshToken}).
		SetResult(TokenResp{}).
		SetError(ErrResp{}).
		Post(BaseURL + "/v2/account/token")
	if err != nil {
		return nil, err
	}

	if e, ok := resp.Error().(*ErrResp); ok {
		return nil, fmt.Errorf("token: %s, %s", e.Code, e.Message)
	}

	res := resp.Result().(*TokenResp)

	err1 := mc.Set(AccessTokenPrefix+refreshToken, res, time.Duration(res.ExpiresIn)*time.Second-time.Minute)
	if err != nil {
		return nil, err1
	}
	err2 := mc.SetDefault(RefreshTokenPrefix+refreshToken, res.RefreshToken)
	if err2 != nil {
		return nil, err2
	}

	return res, nil
}
