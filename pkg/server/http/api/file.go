package api

import (
	"io"
	"net/http"
	"path/filepath"
	"time"

	"turbo-demo/pkg/invoker"
	"turbo-demo/pkg/models/dto"
	"turbo-demo/pkg/utils"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/core/econf"
)

func GetFiles(c *gin.Context) {
	files, err := invoker.FileService.GetFilesList(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "GetFilesList failed: " + err.Error()})
		return
	}
	c.JSON(200, files)
}

func GetFile(c *gin.Context) {
	fileId := c.Param("guid")
	if fileId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get fileId"})
		return
	}
	file, err := invoker.FileService.GetFile(c, fileId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file: " + err.Error()})
		return
	}
	c.JSON(200, file)
}

func DeleteFile(c *gin.Context) {
	fileId := c.Param("guid")
	if fileId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get fileId"})
		return
	}
	err := invoker.FileService.DeleteFile(c, fileId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file: " + err.Error()})
		return
	}
	c.JSON(204, nil)
}

func UploadFile(c *gin.Context) {
	_file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to get file from form"})
		return
	}
	file, err := _file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "file open failed"})
		return
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "file read failed"})
	}
	mimeType := mimetype.Detect(content).String()
	fileName := _file.Filename
	ext := filepath.Ext(fileName)
	f := dto.FileMeta{
		Name:       fileName,
		Size:       int64(len(content)),
		ID:         utils.GenFileGuid(),
		Type:       mimeType,
		CreateTime: time.Now().Unix(),
		Ext:        ext,
	}
	err = invoker.FileService.UploadFileMeta(c, f, content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "file upload failed" + err.Error()})
	}
	c.JSON(200, f)
}

func GetPageParams(c *gin.Context) {
	fileId := c.Param("guid")
	if fileId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get fileId"})
		return
	}
	file, err := invoker.FileService.GetFile(c, fileId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file: " + err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"file":     file,
		"endpoint": econf.GetString("host.previewUrlPrefix"),
		"token":    utils.SignJWT(1, 0),
	})
}
