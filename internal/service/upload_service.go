package service

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"blog/pkg/config"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/google/uuid"
)

// UploadService 文件上传服务接口
type UploadService interface {
	Upload(ctx context.Context, file io.Reader, filename string) (string, error)
}

type uploadService struct {
	client *oss.Client
	cfg    config.OSSConfig
}

// NewUploadService 创建上传服务
func NewUploadService(cfg config.OSSConfig) (UploadService, error) {
	if cfg.AccessKeyID == "" || cfg.AccessKeySecret == "" || cfg.Bucket == "" || cfg.Endpoint == "" {
		return nil, fmt.Errorf("oss 配置不完整")
	}

	ossCfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.AccessKeySecret, "")).
		WithRegion(cfg.Endpoint)

	return &uploadService{
		client: oss.NewClient(ossCfg),
		cfg:    cfg,
	}, nil
}

// Upload 上传文件到 OSS，返回访问 URL
func (s *uploadService) Upload(ctx context.Context, file io.Reader, filename string) (string, error) {
	// 生成唯一文件名：images/2026/08/uuid.ext
	ext := strings.ToLower(filepath.Ext(filename))
	objectKey := fmt.Sprintf("images/%s/%s%s", time.Now().Format("2006/01"), uuid.New().String(), ext)

	request := &oss.PutObjectRequest{
		Bucket: oss.Ptr(s.cfg.Bucket),
		Key:    oss.Ptr(objectKey),
		Body:   file,
	}

	if _, err := s.client.PutObject(ctx, request); err != nil {
		return "", fmt.Errorf("上传文件失败: %w", err)
	}

	// 返回访问 URL
	url := fmt.Sprintf("https://%s.oss-%s.aliyuncs.com/%s", s.cfg.Bucket, s.cfg.Endpoint, objectKey)
	return url, nil
}
