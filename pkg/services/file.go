package services

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"turbo-demo/consts"
	"turbo-demo/pkg/models/dto"
	"turbo-demo/pkg/models/leveldb"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
)

type FileService struct {
	db         *leveldb.LevelDB
	awsService *AwsService
}

func NewFileService(db *leveldb.LevelDB, awsService *AwsService) *FileService {
	return &FileService{
		db:         db,
		awsService: awsService,
	}
}

func (f *FileService) InitCaseFile() {
	// 遍历目录中的所有文件
	err := filepath.Walk(econf.GetString("case.resourcePath"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 忽略目录，只处理文件
		if info.IsDir() {
			return nil
		}
		// 读取文件内容
		data, err := os.ReadFile(path)
		if err != nil {
			log.Printf("无法读取文件 %s: %v", path, err)
			return nil
		}
		fileName := filepath.Base(path)

		// 使用相对路径作为键，将文件内容作为值存入LevelDB
		key := fmt.Sprintf("case_%s", fileName)
		ext := filepath.Ext(fileName)
		err = f.UploadFileMeta(&gin.Context{}, dto.FileMeta{
			ID:         key,
			Name:       fileName,
			CreateTime: time.Now().Unix(),
			Size:       info.Size(),
			Type:       mimetype.Detect(data).String(),
			Ext:        ext,
		}, data)
		if err != nil {
			elog.Panic(fmt.Sprintf("case文件存储失败: %v", err))
		}
		fmt.Printf("文件 %s 已存储到对象存储。\n", path)
		return nil
	})

	// 如果遍历过程中发生错误
	if err != nil {
		elog.Panic(fmt.Sprintf("遍历目录时发生错误: %v", err))
	}
}

func (f *FileService) GetFilesList(c *gin.Context) (files []dto.FileMeta, err error) {
	iter := f.db.DB.NewIterator(nil, nil)
	defer iter.Release()
	for iter.Next() {
		key := string(iter.Key())
		if key == "" || strings.HasPrefix(key, "convert_") || strings.HasPrefix(key, "case_") || strings.HasPrefix(key, "content_") {
			continue
		}
		metaData, err := f.db.GetFileMeta(key)
		if err != nil {
			continue // 跳过无法访问的文件
		}
		files = append(files, metaData)
	}

	if err = iter.Error(); err != nil {
		return nil, err
	}
	return files, nil
}

func (f *FileService) UploadFileMeta(c *gin.Context, file dto.FileMeta, content []byte) error {
	// 存储文件元数据
	err := f.db.SetFileMeta(file.ID, file)
	if err != nil {
		return err
	}
	uploadUrl, err := f.awsService.GetUploadURL(file.Name, file.ID, 3600)
	if err != nil {
		return err
	}
	// 调用下载链接存储到存储空间
	req, _ := http.NewRequest("PUT", uploadUrl, bytes.NewReader(content))
	_, err = http.DefaultClient.Do(req)
	return err
}

func (f *FileService) GetUploadUrl(c *gin.Context, path, fileId string) (string, error) {
	return f.awsService.GetUploadURL(path, fileId, 3600)
}

func (f *FileService) GetFile(c *gin.Context, fileId string) (file dto.FileMeta, err error) {
	return f.db.GetFileMeta(fileId)
}

func (f *FileService) GetDownloadUrl(fileId, fileName, path, disposition string, exp int64) (url string, err error) {
	if exp == 0 {
		exp = 3600
	}
	downloadUrl, err := f.awsService.GetDownloadURL(fileId, path, fileName, disposition, exp)
	if err != nil {
		return "", err
	}
	return downloadUrl, nil
}

func (f *FileService) DeleteFile(c *gin.Context, fileId string) (err error) {
	err = f.db.DeleteFileMeta(fileId)
	if err != nil {
		return err
	}
	return f.awsService.DeleteFolder(fileId + "/")
}

func (f *FileService) CheckContentExist(c *gin.Context, fileId string) bool {
	return f.awsService.FileExists(fmt.Sprintf("%s/%s/%s", consts.Dir, fileId, "content"))
}
