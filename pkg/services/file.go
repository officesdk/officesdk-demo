package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"turbo-demo/pkg/models/dto"
	"turbo-demo/pkg/models/leveldb"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
)

type FileService struct {
	db *leveldb.LevelDB
}

func NewFileService(db *leveldb.LevelDB) *FileService {
	return &FileService{
		db: db,
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
		err = f.db.SetFileMeta(key, dto.FileMeta{
			ID:         key,
			Name:       fileName,
			CreateTime: time.Now().Unix(),
			Size:       info.Size(),
			Type:       mimetype.Detect(data).String(),
		})
		if err != nil {
			elog.Panic(fmt.Sprintf("存储文件元信息 %s 到LevelDB 失败: %v", path, err))
		}
		fmt.Printf("文件 %s 已存储到LevelDB。\n", path)
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

func (f *FileService) UploadFile(c *gin.Context, file dto.FileMeta, content []byte) error {
	// 存储文件元数据
	err := f.db.SetFileMeta(file.ID, file)
	if err != nil {
		return err
	}
	// 存储文件内容
	return f.WriteBytesToFile(content, UploadFilePath(file.ID, file.Ext))
}

func (f *FileService) GetFile(c *gin.Context, fileId string) (file dto.FileMeta, err error) {
	return f.db.GetFileMeta(fileId)
}

func (f *FileService) GetDownloadUrl(fileId string) (url string) {
	host := econf.GetString("host.downloadUrlPrefix")
	return fmt.Sprintf("%s/showcase/%s/download", host, fileId)
}

func (f *FileService) DeleteFile(c *gin.Context, fileId string) (err error) {
	err = f.db.DeleteFileMeta(fileId)
	if err != nil {
		return err
	}
	return f.DeleteFileContent(fmt.Sprintf("%s/%s", econf.GetString("case.filepath"), fileId))
}

// WriteBytesToFile 将字节数据写入指定文件，如果目录或文件不存在则创建
func (f *FileService) WriteBytesToFile(data []byte, filePath string) error {
	// 获取文件所在目录
	dir := filepath.Dir(filePath)

	// 检查目录是否存在，如果不存在则创建
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755) // 创建目录，权限为 0755
		if err != nil {
			return err
		}
	}

	// 写入文件，如果文件不存在则创建，权限为 0644
	err := os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileService) DeleteFileContent(filePath string) (err error) {
	// 检查路径是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // 如果文件不存在，直接返回成功
	}

	// 尝试删除文件或目录
	err = os.RemoveAll(filePath)
	if err != nil {
		return fmt.Errorf("删除文件失败: %v", err)
	}
	return nil
}

func (f *FileService) GetFileContent(fileId string) (content []byte, err error) {
	file, err := f.db.GetFileMeta(fileId)
	if err != nil {
		return nil, err
	}
	filePath := ""
	if strings.HasPrefix(fileId, "case_") {
		filePath = ResourceFilePath(fileId[5:], file.Ext)
	} else {
		filePath = UploadFilePath(fileId, file.Ext)
	}
	return os.ReadFile(filePath)
}

func UploadFilePath(fileID string, fileExt string) string {
	return filepath.Join(econf.GetString("case.filepath"), fileID, fmt.Sprintf("source%s", fileExt))
}

func ResourceFilePath(fileID string, fileExt string) string {
	return filepath.Join(econf.GetString("case.resourcePath"), fileID, fileExt)
}
