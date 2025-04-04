package callback

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/core/econf"

	"turbo-demo/pkg/invoker"
)

func Verify(c *gin.Context) {
	// fileId := c.Param("fileId")
	// invoker.Leveldb.GetFile(fileId)
	c.JSON(200, gin.H{
		"code": 0,
		"data": gin.H{
			"currentUserInfo": gin.H{
				"id":     "123456",
				"name":   "Void",
				"avatar": "",
				"email":  "",
			},
		},
	})
}
func GetFile(c *gin.Context) {
	fileId := c.Param("fileId")
	file, err := invoker.Leveldb.GetFile(fileId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "0",
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"data": file,
	})
}

func GetFileDownload(c *gin.Context) {
	fileId := c.Param("fileId")
	c.JSON(200, gin.H{
		"code": 0,
		"data": gin.H{
			"url": invoker.FileService.GetDownloadUrl(fileId),
		},
	})
}

func GetFileWatermark(c *gin.Context) {
	fileId := c.Param("fileId")
	// todo 暂时无水印设置功能
	c.JSON(200, gin.H{
		"code":    0,
		"message": "",
		"data": gin.H{
			"fillStyle":  "rgba( 192, 192, 192, 0.6 )",
			"font":       "bold 20px Serif",
			"horizontal": 50,
			"rotate":     -0.7853982,
			"type":       1,
			"value":      fmt.Sprintf("%s\n%s", fileId, time.Now().In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05")),
			"vertical":   100,
		},
	})
}

// UploadAddress 获取上传地址 todo
func UploadAddress(c *gin.Context) {
	fileId := c.Param("fileId")
	c.JSON(200, gin.H{
		"code":    0,
		"message": "",
		"data": gin.H{
			"url":          fmt.Sprintf("%s/v1/callback/files/%s/upload", econf.GetString("host.downloadUrlPrefix"), fileId),
			"fileFieldKey": "file",
			"params": gin.H{
				"AccessKeyId": "",
			},
		},
	})
}

func UploadFile(c *gin.Context) {
	// fileId := c.Param("fileId")
	// _file, err := c.FormFile("file")
	// if err != nil {
	// 	c.JSON(400, gin.H{"error": "Failed to get file from form"})
	// 	return
	// }
	// file, err := _file.Open()
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": "file open failed"})
	// 	return
	// }
	// defer file.Close()
	// content, err := io.ReadAll(file)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": "file read failed"})
	// }
	// mimeType := mimetype.Detect(content).String()
	// fileName := _file.Filename
	// f := dto.File{
	// 	Name:       fileName,
	// 	Size:       int64(len(content)),
	// 	ID:         fmt.Sprintf("convert_%s", fileId),
	// 	Type:       mimeType,
	// 	CreateTime: time.Now().Unix(),
	// 	Content:    content,
	// }
	// err = invoker.FileService.UploadFile(c, f)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": "file upload failed" + err.Error()})
	// }
	c.JSON(200, gin.H{
		"code": 0,
		"data": gin.H{},
	})
}

// UploadComplete 上传完成后, 回调通知上传结果 todo
func UploadComplete(c *gin.Context) {
	fileId := c.Param("fileId")
	file, _ := invoker.FileService.GetFile(c, fileId)
	c.JSON(200, gin.H{
		"code": 0,
		"data": file,
	})
}

// OpenAIConfig 定义OpenAI配置结构
type OpenAIConfig struct {
	AIIcon string      `toml:"aiIcon"`
	LLM    []LLMConfig `toml:"llm"`
}

// LLMConfig 定义LLM配置结构
type LLMConfig struct {
	BaseUrl        string `toml:"baseUrl"`
	TextModel      string `toml:"textModel"`
	Token          string `toml:"token"`
	Name           string `toml:"name"`
	ProxyUrl       string `toml:"proxyUrl"`
	Subservice     string `toml:"subservice"`
	InputMaxToken  int    `toml:"inputMaxToken"`
	OutputMaxToken int    `toml:"outputMaxToken"`
}

func AIConfig(c *gin.Context) {
	openaiConfig := OpenAIConfig{}
	econf.UnmarshalKey("openai", &openaiConfig)

	c.JSON(200, gin.H{
		"code": 0,
		"data": gin.H{
			"aiIcon":  openaiConfig.AIIcon,
			"llmList": openaiConfig.LLM,
		},
	})
}
