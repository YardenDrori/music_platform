package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/minio/minio-go/v7"

	"github.com/YardenDrori/music-platform/internal/apperrors"
)

type service struct {
	s3Core               *minio.Core
	preginedURLsDuration time.Duration
}

func NewService(client *minio.Client, presignedURLsDuration time.Duration) Service {
	return &service{
		s3Core:               &minio.Core{Client: client},
		preginedURLsDuration: presignedURLsDuration,
	}
}

func (s *service) PresignedUpload(
	ctx context.Context,
	bucketName string,
	objectKey string,
) (*url.URL, error) {
	uploadUrl, err := s.s3Core.PresignedPutObject(
		ctx,
		bucketName,
		objectKey,
		s.preginedURLsDuration,
	)
	if err != nil {
		return nil, fmt.Errorf("generating a presigned url for a put request: %w", err)
	}
	return uploadUrl, nil
}

func (s *service) PutObject(
	ctx context.Context,
	bucketName string,
	objectKey string,
	reader io.Reader,
	size int64,
	opts PutOptions,
) error {
	minioOpts := minio.PutObjectOptions{
		ContentType:    opts.ContentType,
		SendContentMd5: opts.SendContentMD5,
	}
	_, err := s.s3Core.Client.PutObject(ctx, bucketName, objectKey, reader, size, minioOpts)
	if err != nil {
		return fmt.Errorf(
			"sending put request to minio storage: %w",
			apperrors.NewErrInternal().WithCause(err),
		)
	}
	return nil
}

func (s *service) InitiateMultipartUpload(
	ctx context.Context,
	bucketName string,
	objectKey string,
	opts PutOptions,
) (string, error) {
	minioOpts := minio.PutObjectOptions{
		ContentType:    opts.ContentType,
		SendContentMd5: opts.SendContentMD5,
	}

	uploadID, err := s.s3Core.NewMultipartUpload(ctx, bucketName, objectKey, minioOpts)
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
			s.preginedURLsDuration,
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
	opts PutOptions,
) error {

	minioOpts := minio.PutObjectOptions{
		ContentType:    opts.ContentType,
		SendContentMd5: opts.SendContentMD5,
	}

	_, err := s.s3Core.CompleteMultipartUpload(
		ctx,
		bucketName,
		objectKey,
		uploadID,
		ETags,
		minioOpts,
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
