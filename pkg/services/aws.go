package services

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"turbo-demo/consts"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gotomicro/cetus/l"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
)

type AwsService struct {
	client *s3.S3
	bucket string
}

func NewAwsService() (*AwsService, error) {
	// 从配置文件读取配置
	accessKeyId := econf.GetString("awos.accessKeyId")
	accessKeySecret := econf.GetString("awos.accessKeySecret")
	endpoint := econf.GetString("awos.endpoint")
	region := econf.GetString("awos.region")
	bucket := econf.GetString("awos.bucket")

	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKeyId, accessKeySecret, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	return &AwsService{
		client: s3.New(sess),
		bucket: bucket,
	}, nil
}

// GetUploadURL 生成上传URL，有效期为指定的秒数
func (s *AwsService) GetUploadURL(path, key string, expireSeconds int64) (string, error) {
	if expireSeconds < 86400 { // 默认最小值为一天
		expireSeconds = 86400
	}
	req, _ := s.client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fmt.Sprintf("%s/%s/%s", consts.Dir, key, path)),
	})

	return req.Presign(time.Duration(expireSeconds) * time.Second)
}

// GetDownloadURL 生成下载URL，有效期为指定的秒数
// disposition: "inline" - 在浏览器中直接显示；"attachment" - 作为附件下载
func (s *AwsService) GetDownloadURL(key, path, filename, disposition string, expireSeconds int64) (string, error) {
	encodedFilename := url.PathEscape(filename)

	// 根据 disposition 设置 Content-Disposition
	var contentDisposition string
	if disposition == "inline" {
		contentDisposition = fmt.Sprintf(`inline; filename="%s"`, encodedFilename)
	} else {
		contentDisposition = fmt.Sprintf(`attachment; filename="%s"`, encodedFilename)
	}

	// 创建请求对象
	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket:                     aws.String(s.bucket),
		Key:                        aws.String(fmt.Sprintf("%s/%s/%s", consts.Dir, key, path)),
		ResponseContentDisposition: aws.String(contentDisposition),
	})

	// 生成预签名的 URL，添加必要的参数
	u, err := req.Presign(time.Duration(expireSeconds) * time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %v", err)
	}

	return u, nil
}

// DeleteFile 从对象存储删除文件
func (s *AwsService) DeleteFile(key string) error {
	_, err := s.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fmt.Sprintf("%s/%s", consts.Dir, key)),
	})
	return err
}

// DeleteFolder 删除文件夹及其中的所有文件
func (s *AwsService) DeleteFolder(folderKey string) error {
	// 列出文件夹内的所有对象
	listObjectsInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(fmt.Sprintf("%s/%s", consts.Dir, folderKey)),
	}

	// 获取对象列表
	resp, err := s.client.ListObjectsV2(listObjectsInput)
	if err != nil {
		return fmt.Errorf("无法列出对象: %v", err)
	}

	// 如果没有文件，直接返回
	if len(resp.Contents) == 0 {
		return nil
	}

	// 删除文件夹内的所有文件
	for _, obj := range resp.Contents {
		_, err := s.client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(*obj.Key),
		})
		if err != nil {
			elog.Error("Delete Error", l.E(err))
		}
	}

	// 返回成功删除信息
	return nil
}

// FileExists 检查指定的文件是否存在于 S3 中
func (s *AwsService) FileExists(key string) bool {
	// 使用 HeadObject 请求检查文件是否存在
	_, err := s.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	// 判断错误类型
	if err == nil {
		// 文件存在
		return true
	}
	// 如果是文件不存在的错误，返回 false
	var aerr awserr.Error
	if errors.As(err, &aerr) {
		return aerr.Code() == s3.ErrCodeNoSuchKey
	}
	elog.Error("Failed to check file existence: ", l.E(err))
	return false
}
