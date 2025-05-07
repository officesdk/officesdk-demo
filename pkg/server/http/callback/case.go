package callback

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"turbo-demo/pkg/invoker"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/cetus/l"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"github.com/officesdk/go-sdk/officesdk"
)

type FileProvider struct{}

func (f *FileProvider) VerifyFile(c *gin.Context, fileId string) (*officesdk.VerifyResponse, error) {
	_userId, _ := c.Get("userId")
	userId := fmt.Sprintf("%d", _userId)
	return &officesdk.VerifyResponse{
		CurrentUserInfo: officesdk.UserInfo{
			ID:    userId,
			Name:  "demo",
			Email: "a@b.com",
		},
	}, nil
}

func (f *FileProvider) GetFile(c *gin.Context, fileId string) (*officesdk.FileResponse, error) {
	file, err := invoker.Leveldb.GetFileMeta(fileId)
	if err != nil {
		return nil, err
	}
	fromSDK := invoker.FileService.CheckContentExist(c, fileId)
	return &officesdk.FileResponse{
		ID:         file.ID,
		Name:       file.Name,
		Version:    uint32(file.Version),
		CreateTime: file.CreateTime,
		ModifyTime: file.ModifyTime,
		CreatorID:  file.CreatorId,
		ModifierID: file.ModifierId,
		FromSDK:    fromSDK,
	}, nil
}

func (f *FileProvider) GetFileDownload(c *gin.Context, fileId string) (*officesdk.DownloadResponse, error) {
	file, err := invoker.Leveldb.GetFileMeta(fileId)
	if err != nil {
		return nil, err
	}
	downloadUrl, err := invoker.FileService.GetDownloadUrl(fileId, file.Name, file.Name, "", 0)
	if err != nil {
		return nil, err
	}
	elog.Info("GetFileDownload", l.S("downloadUrl", downloadUrl))

	return &officesdk.DownloadResponse{
		URL: downloadUrl,
	}, nil
}

func (f *FileProvider) GetFileWatermark(c *gin.Context, fileId string) (*officesdk.WatermarkResponse, error) {
	// todo 暂时无水印设置功能
	return &officesdk.WatermarkResponse{
		Type:       1,
		Value:      fmt.Sprintf("%s\n%s", fileId, time.Now().In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05")),
		FillStyle:  "rgba( 192, 192, 192, 0.6 )",
		Font:       "bold 20px Serif",
		Rotate:     -0.7853982,
		Horizontal: 50,
		Vertical:   100,
	}, nil
}

type UploadBody struct {
	ObjectName  string `json:"object_name"`
	ContentType string `json:"content_type"`
}

// GetUploadURL 上传文件转码信息
func (f *FileProvider) GetUploadURL(c *gin.Context, fileId string) (*officesdk.UploadURLResponse, error) {
	body := UploadBody{}
	err := c.BindJSON(&body)
	if err != nil || body.ObjectName == "" {
		elog.Error("GetUploadURL body err: ", l.E(err))
		return nil, errors.New("body is required")
	}

	url, err := invoker.FileService.GetUploadUrl(c, body.ObjectName, fileId)
	if err != nil {
		return nil, errors.New("GetUploadUrl failed")
	}

	return &officesdk.UploadURLResponse{
		URL:    url,
		Method: "PUT",
		Headers: map[string]string{
			"Content-Type": body.ContentType,
		},
		Params: map[string]string{
			"X-Test-Token": "test",
		},
		CompletionParams: map[string]string{
			"test1": "test1",
			"test2": "test2",
			"test3": "test3",
		},
	}, nil
}

// CompleteUpload 上传文件转码完成
func (f *FileProvider) CompleteUpload(c *gin.Context, fileId string) (*officesdk.UploadCompletionResponse, error) {
	file, _ := invoker.Leveldb.GetFileMeta(fileId)
	// 打印请求参数
	elog.Info("CompleteUpload request params",
		l.S("file_id", fileId),
		l.S("object_name", c.Query("object_name")),
		l.S("object_size", c.Query("object_size")),
		l.S("content_type", c.Query("content_type")),
	)
	// 读取请求体
	var requestBody map[string]interface{}
	if err := c.ShouldBindJSON(&requestBody); err == nil {
		elog.Info("CompleteUpload request body",
			l.A("request", requestBody["request"]),
			l.A("digest", requestBody["digest"]),
			l.A("completion_params", requestBody["completion_params"]),
		)
	}
	// 打印响应参数
	elog.Info("CompleteUpload response",
		l.A("response", requestBody["response"]),
		l.S("status_code", c.Query("status_code")),
		l.A("headers", requestBody["headers"]),
		l.A("body", requestBody["body"]),
	)
	return &officesdk.UploadCompletionResponse{
		ID:         file.ID,
		Version:    int(file.Version),
		CreateTime: file.CreateTime,
		ModifyTime: file.ModifyTime,
		CreatorID:  file.CreatorId,
		ModifierID: file.ModifierId,
	}, nil
}

// GetDownloadURL 下载文件转码信息
func (f *FileProvider) GetDownloadURL(c *gin.Context, fileId string) (*officesdk.DownloadResponse, error) {
	objName := c.Query("object_name")
	if objName == "" {
		return nil, errors.New("object_name is required")
	}
	expS := c.Query("expires_in")
	disposition := c.Query("disposition")
	var exp int64
	var err error
	if expS != "" {
		exp, err = strconv.ParseInt(expS, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	url, err := invoker.FileService.GetDownloadUrl(fileId, objName, objName, disposition, exp)
	elog.Info("GetDownloadUrl", l.S("name", objName), l.S("url", url))
	if err != nil {
		return nil, err
	}
	return &officesdk.DownloadResponse{
		URL: url,
	}, nil
}

// GetAssetUploadURL 上传文件附件资源信息
func (f *FileProvider) GetAssetUploadURL(c *gin.Context, fileId string) (*officesdk.AssetUploadURLResponse, error) {
	body := UploadBody{}
	err := c.BindJSON(&body)
	if err != nil || body.ObjectName == "" {
		elog.Error("GetAssetUploadURL body err: ", l.E(err))
		return nil, errors.New("object_name is required")
	}

	url, err := invoker.FileService.GetUploadUrl(c, body.ObjectName, fileId)
	if err != nil {
		return nil, errors.New("GetUploadUrl failed")
	}

	return &officesdk.AssetUploadURLResponse{
		URL:    url,
		Method: "PUT",
		Headers: map[string]string{
			"Content-Type": body.ContentType,
		},
		Params: map[string]string{
			"X-Test-Token": "test",
		},
		CompletionParams: map[string]string{
			"test1": "test1",
			"test2": "test2",
			"test3": "test3",
		},
	}, nil
}

// AssetCompleteUpload 上传文件附件资源完成
func (f *FileProvider) AssetCompleteUpload(c *gin.Context, fileId string) (*officesdk.UploadCompletionResponse, error) {
	file, _ := invoker.Leveldb.GetFileMeta(fileId)
	// 打印请求参数
	elog.Info("AssetCompleteUpload request params",
		l.S("file_id", fileId),
		l.S("object_name", c.Query("object_name")),
		l.S("object_size", c.Query("object_size")),
		l.S("content_type", c.Query("content_type")),
	)
	// 读取请求体
	var requestBody map[string]interface{}
	if err := c.ShouldBindJSON(&requestBody); err == nil {
		elog.Info("AssetCompleteUpload request body",
			l.A("request", requestBody["request"]),
			l.A("digest", requestBody["digest"]),
			l.A("completion_params", requestBody["completion_params"]),
		)
	}
	// 打印响应参数
	elog.Info("AssetCompleteUpload response",
		l.A("response", requestBody["response"]),
		l.S("status_code", c.Query("status_code")),
		l.A("headers", requestBody["headers"]),
		l.A("body", requestBody["body"]),
	)
	return &officesdk.UploadCompletionResponse{
		ID:         file.ID,
		Version:    int(file.Version),
		CreateTime: file.CreateTime,
		ModifyTime: file.ModifyTime,
		CreatorID:  file.CreatorId,
		ModifierID: file.ModifierId,
	}, nil
}

// GetAssetDownloadURL 下载文件附件资源信息
func (f *FileProvider) GetAssetDownloadURL(c *gin.Context, fileId string) (*officesdk.DownloadResponse, error) {
	objName := c.Query("object_name")
	if objName == "" {
		return nil, errors.New("object_name is required")
	}
	expS := c.Query("expires_in")
	disposition := c.Query("disposition")
	var exp int64
	var err error
	if expS != "" {
		exp, err = strconv.ParseInt(expS, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	url, err := invoker.FileService.GetDownloadUrl(fileId, objName, objName, disposition, exp)
	elog.Info("GetAssetDownloadURL", l.S("name", objName), l.S("url", url))
	if err != nil {
		return nil, err
	}
	return &officesdk.DownloadResponse{
		URL: url,
	}, nil
}

// AIProvider 实现 AI 相关接口
type AIProvider struct{}

func (p *AIProvider) AIConfig(c *gin.Context) (*officesdk.AIConfigResponse, error) {
	openaiConfig := officesdk.AIConfigResponse{}
	err := econf.UnmarshalKey("openai", &openaiConfig)
	if err != nil {
		return nil, err
	}
	return &openaiConfig, nil
}
