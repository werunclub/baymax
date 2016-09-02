package main

import (
	"io"

	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pborman/uuid"
)

func getBucket() (*oss.Bucket, error) {
	client, err := oss.New(
		Config.Storage.Endpoint,
		Config.Storage.AccessKey,
		Config.Storage.AccessSecret)

	if err != nil {
		return nil, err
	}

	bucket, err := client.Bucket(Config.Storage.BucketName)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

func newFileName() string {
	return uuid.NewUUID().String()
}

func getUrl(fileKey string) string {
	return fmt.Sprintf("%v/%v", Config.Url.BaseUrl, fileKey)
}

func storePicture(objectKey string, reader io.Reader, contentType string) error {
	bucket, err := getBucket()
	if err != nil {
		return err
	}

	return bucket.PutObject(objectKey, reader, oss.ContentType(contentType))
}
