package storage

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/minio/minio-go/v7"
)

type service struct {
	s3Core *minio.Core
}

func NewService(client *minio.Client) Service {
	return &service{s3Core: &minio.Core{Client: client}}
}

func (s *service) InitiateMultipartUpload(
	ctx context.Context,
	bucketName string,
	objectKey string,
	opts minio.PutObjectOptions,
) (string, error) {
	uploadID, err := s.s3Core.NewMultipartUpload(ctx, bucketName, objectKey, opts)
	if err != nil {
		return "", fmt.Errorf("initiating new presigned multipart upload: %w", err)
	}
	return uploadID, nil
}

func (s *service) GetPresignedMultipartPartsURLs(
	ctx context.Context,
	bucketName string,
	objectKey string,
	uploadID string,
	totalPartsCount int,
) ([]string, error) {
	var presignedURLs []string
	for idx := range totalPartsCount {
		req := url.Values{
			"partNumber": {strconv.Itoa(idx + 1)},
			"uploadId":   {uploadID},
		}
		url, err := s.s3Core.Presign(
			ctx,
			http.MethodPut,
			bucketName,
			objectKey,
			time.Duration(15)*time.Minute,
			req,
		)
		if err != nil {
			return nil, fmt.Errorf("generating presigned multipart upload URLs: %w", err)
		}
		presignedURLs = append(presignedURLs, url.String())
	}

	return presignedURLs, nil
}

func (s *service) CompleteMultipartUpload(
	ctx context.Context,
	bucketName string,
	objectKey string,
	uploadID string,
	ETags []minio.CompletePart,
	opts minio.PutObjectOptions,
) error {
	_, err := s.s3Core.CompleteMultipartUpload(
		ctx,
		bucketName,
		objectKey,
		uploadID,
		ETags,
		opts,
	)
	if err != nil {
		anotherErr := s.AbortMultipartUpload(context.Background(), bucketName, objectKey, uploadID)
		if anotherErr != nil {
			return fmt.Errorf(
				"completing presgined multipart upload: %w %w",
				err,
				anotherErr,
			)
		}
		return fmt.Errorf("completing presgined multipart upload: %w", err)
	}

	return nil
}

func (s *service) AbortMultipartUpload(
	ctx context.Context,
	bucketName string,
	objectKey string,
	uploadID string,
) error {
	err := s.s3Core.AbortMultipartUpload(ctx, bucketName, objectKey, uploadID)
	if err != nil {
		return fmt.Errorf("aborting multipart upload %w", err)
	}
	return nil
}
