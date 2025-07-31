package oss

import (
	"fmt"
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OSSClient struct {
	Endpoint string
	Client   *oss.Client
	Bucket   *oss.Bucket
}

func NewOSSClient(endpoint, accessKeyID, accessKeySecret, bucketName string) (*OSSClient, error) {
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("创建 OSS 客户端失败: %w", err)
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, fmt.Errorf("获取 Bucket 失败: %w", err)
	}
	return &OSSClient{
		Endpoint: endpoint,
		Client:   client,
		Bucket:   bucket,
	}, nil
}

func (c *OSSClient) Upload(localFile, objectKey string) (string, error) {
	err := c.Bucket.PutObjectFromFile(objectKey, localFile, oss.ObjectACL(oss.ACLPublicRead))
	if err != nil {
		return "", fmt.Errorf("上传文件失败: %w", err)
	}
	url := fmt.Sprintf("https://%s.%s/%s", c.Bucket.BucketName, c.Endpoint, objectKey)
	return url, nil
}

func (c *OSSClient) UploadReader(objectKey string, reader io.Reader) (string, error) {
	bucket, err := c.Client.Bucket(c.Bucket.BucketName)
	if err != nil {
		return "", fmt.Errorf("获取 Bucket 失败: %w", err)
	}
	err = bucket.PutObject(objectKey, reader, oss.ObjectACL(oss.ACLPublicRead))
	if err != nil {
		return "", fmt.Errorf("上传文件失败: %w", err)
	}
	url := fmt.Sprintf("https://%s.%s/%s", c.Bucket.BucketName, c.Endpoint, objectKey)
	return url, nil
}
