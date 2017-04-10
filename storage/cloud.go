package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"errors"
	"bytes"
	"io"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/MapOnline/go-core/logger"
)

type CloudStorageConfiguration struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	Region string `json:"region"`
	Bucket string `json:"bucket"`
}

type CloudStorage struct {
	Configuration *CloudStorageConfiguration
}

func NewCloudStorage(Configuration *CloudStorageConfiguration)(StorageHandle) {
	return StorageHandle(&CloudStorage{
		Configuration: Configuration,
	})
}

func (Storage *CloudStorage) Service()(*s3.S3, error) {
	var sess, err = session.NewSession(&aws.Config{
		Region: aws.String(Storage.Configuration.Region),
		Credentials: credentials.NewStaticCredentials(Storage.Configuration.AccessKey, Storage.Configuration.SecretKey, ""),
	})

	if err != nil {
		return nil, err
	}

	return s3.New(sess), nil
}

func (Storage *CloudStorage) HasFile(file string) bool {
	var err error
	svc, err := Storage.Service()

	if err != nil {
		logger.Log.Println("Cannot connect to S3")
		return false
	}

	params := &s3.HeadObjectInput{
		Bucket: aws.String(Storage.Configuration.Bucket),
		Key: aws.String(file),
	}

	_, err = svc.HeadObject(params)
	return err == nil
}

func (Storage *CloudStorage) GetFile(file string)(io.ReadCloser, error) {
	svc, err := Storage.Service()

	if err != nil {
		return nil, err
	}

	params := &s3.GetObjectInput{
		Bucket: aws.String(Storage.Configuration.Bucket),
		Key: aws.String(file),
	}

	resp, err := svc.GetObject(params)

	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (Storage *CloudStorage) DeleteFile(file string) error {
	if !Storage.HasFile(file) {
		return errors.New("File doesn't exist")
	}

	svc, err := Storage.Service()

	if err != nil {
		return err
	}

	params := &s3.DeleteObjectInput{
		Bucket: aws.String(Storage.Configuration.Bucket),
		Key: aws.String(file),
	}

	_, err = svc.DeleteObject(params)
	return err
}

func (Storage *CloudStorage) WriteFile(buffer []byte, file string) error {
	if Storage.HasFile(file) {
		return errors.New("File already exists")
	}

	svc, err := Storage.Service()

	if err != nil {
		return err
	}

	params := &s3.PutObjectInput{
		Key: aws.String(file),
		Bucket: aws.String(Storage.Configuration.Bucket),
		Body: bytes.NewReader(buffer),
	}

	_, err = svc.PutObject(params)
	return err
}