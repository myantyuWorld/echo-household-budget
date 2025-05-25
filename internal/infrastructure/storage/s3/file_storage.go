package s3

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"echo-household-budget/internal/domain/repository"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3FileStorage struct {
	client     *s3.Client
	bucketName string
}

func NewS3FileStorage(client *s3.Client, bucketName string) repository.FileStorageRepository {
	return &S3FileStorage{
		client:     client,
		bucketName: bucketName,
	}
}

func (s *S3FileStorage) UploadFile(fileData []byte, fileName string) (string, error) {
	ctx := context.Background()

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
		Body:   bytes.NewReader(fileData),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return s.GetFileURL(fileName)
}

func (s *S3FileStorage) GetFileURL(fileName string) (string, error) {
	presignClient := s3.NewPresignClient(s.client)
	presignedURL, err := presignClient.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
	}, s3.WithPresignExpires(time.Hour*24)) // URLの有効期限を24時間に設定

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL.URL, nil
}

func (s *S3FileStorage) DeleteFile(fileName string) error {
	ctx := context.Background()

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}
