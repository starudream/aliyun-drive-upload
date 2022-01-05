package aliyunDriveUpload

import (
	"encoding/xml"
	"time"
)

type ErrResp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type TokenReq struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

type TokenResp struct {
	DefaultSboxDriveId string    `json:"default_sbox_drive_id"`
	Role               string    `json:"role"`
	DeviceId           string    `json:"device_id"`
	UserName           string    `json:"user_name"`
	NeedLink           bool      `json:"need_link"`
	ExpireTime         time.Time `json:"expire_time"`
	PinSetup           bool      `json:"pin_setup"`
	NeedRpVerify       bool      `json:"need_rp_verify"`
	Avatar             string    `json:"avatar"`
	TokenType          string    `json:"token_type"`
	AccessToken        string    `json:"access_token"`
	DefaultDriveId     string    `json:"default_drive_id"`
	DomainId           string    `json:"domain_id"`
	RefreshToken       string    `json:"refresh_token"`
	IsFirstLogin       bool      `json:"is_first_login"`
	UserId             string    `json:"user_id"`
	NickName           string    `json:"nick_name"`
	State              string    `json:"state"`
	ExpiresIn          int       `json:"expires_in"`
	Status             string    `json:"status"`
}

type FileReq struct {
	DriveId       string        `json:"drive_id"`
	PartInfoList  []FilePartReq `json:"part_info_list"`
	ParentFileId  string        `json:"parent_file_id"`
	Name          string        `json:"name"`
	Type          string        `json:"type"`
	CheckNameMode string        `json:"check_name_mode"`
	Size          int           `json:"size"`
	PreHash       string        `json:"pre_hash"`
}

type FilePartReq struct {
	PartNumber        int    `json:"part_number"`
	UploadUrl         string `json:"upload_url,omitempty"`
	InternalUploadUrl string `json:"internal_upload_url,omitempty"`
	ContentType       string `json:"content_type,omitempty"`
}

type FileResp struct {
	ParentFileId string        `json:"parent_file_id"`
	PartInfoList []FilePartReq `json:"part_info_list"`
	UploadId     string        `json:"upload_id"`
	RapidUpload  bool          `json:"rapid_upload"`
	Type         string        `json:"type"`
	FileId       string        `json:"file_id"`
	DomainId     string        `json:"domain_id"`
	DriveId      string        `json:"drive_id"`
	FileName     string        `json:"file_name"`
	EncryptMode  string        `json:"encrypt_mode"`
	Location     string        `json:"location"`
}

type CompleteReq struct {
	DriveId  string `json:"drive_id"`
	UploadId string `json:"upload_id"`
	FileId   string `json:"file_id"`
}

type CompleteResp struct {
	DriveId          string    `json:"drive_id"`
	DomainId         string    `json:"domain_id"`
	FileId           string    `json:"file_id"`
	Name             string    `json:"name"`
	Type             string    `json:"type"`
	ContentType      string    `json:"content_type"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	FileExtension    string    `json:"file_extension"`
	Hidden           bool      `json:"hidden"`
	Size             int       `json:"size"`
	Starred          bool      `json:"starred"`
	Status           string    `json:"status"`
	UserMeta         string    `json:"user_meta"`
	UploadId         string    `json:"upload_id"`
	ParentFileId     string    `json:"parent_file_id"`
	Crc64Hash        string    `json:"crc64_hash"`
	ContentHash      string    `json:"content_hash"`
	ContentHashName  string    `json:"content_hash_name"`
	Category         string    `json:"category"`
	EncryptMode      string    `json:"encrypt_mode"`
	CreatorType      string    `json:"creator_type"`
	CreatorId        string    `json:"creator_id"`
	LastModifierType string    `json:"last_modifier_type"`
	LastModifierId   string    `json:"last_modifier_id"`
	RevisionId       string    `json:"revision_id"`
	Location         string    `json:"location"`
}

type OSSResp struct {
	XMLName           xml.Name `xml:"Error"`
	Code              string   `xml:"Code"`
	Message           string   `xml:"Message"`
	RequestId         string   `xml:"RequestId"`
	HostId            string   `xml:"HostId"`
	OSSAccessKeyId    string   `xml:"OSSAccessKeyId"`
	SignatureProvided string   `xml:"SignatureProvided"`
	StringToSign      string   `xml:"StringToSign"`
}

type DownloadReq struct {
	DriveId   string `json:"drive_id"`
	FileId    string `json:"file_id"`
	ExpireSec int    `json:"expire_sec"`
}

type DownloadResp struct {
	Method      string    `json:"method"`
	Url         string    `json:"url"`
	InternalUrl string    `json:"internal_url"`
	Expiration  time.Time `json:"expiration"`
	Size        int       `json:"size"`
	RateLimit   struct {
		PartSpeed int `json:"part_speed"`
		PartSize  int `json:"part_size"`
	} `json:"ratelimit"`
	Crc64Hash       string `json:"crc64_hash"`
	ContentHash     string `json:"content_hash"`
	ContentHashName string `json:"content_hash_name"`
}
