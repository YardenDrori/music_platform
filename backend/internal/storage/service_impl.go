package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"

	"github.com/YardenDrori/music-platform/internal/apperrors"
)

type service struct {
	s3Core               *minio.Core
	preginedURLsDuration time.Duration
	publicUrlPrefix      *string
}

func NewService(
	client *minio.Client,
	presignedURLsDuration time.Duration,
	publicUrlPrefix *string,
) Service {
	return &service{
		s3Core:               &minio.Core{Client: client},
		preginedURLsDuration: presignedURLsDuration,
		publicUrlPrefix:      publicUrlPrefix,
	}
}

func (s *service) BuildPublicGetUrl(bucketName string, objectKey string) *string {
	result := strings.Join([]string{
		*s.publicUrlPrefix, bucketName, objectKey,
	}, "/")
	return &result
}

func toMinioChecksum(c ChecksumAlgo) minio.ChecksumType {
	switch c {
	case ChecksumSha256:
		return minio.ChecksumSHA256
	default:
		return minio.ChecksumNone
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
		ContentType: opts.ContentType,
		Checksum:    toMinioChecksum(opts.Checksum),
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
		ContentType: opts.ContentType,
		Checksum:    toMinioChecksum(opts.Checksum),
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
		ContentType: opts.ContentType,
		Checksum:    toMinioChecksum(opts.Checksum),
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

func (s *service) DeleteObject(
	ctx context.Context,
	bucketName string,
	objectKey string,
	opts DeleteOptions,
) error {
	minioOpts := minio.RemoveObjectOptions{
		ForceDelete:      opts.ForceDelete,
		GovernanceBypass: opts.GovernanceBypass,
		VersionID:        opts.VersionID,
	}
	if err := s.s3Core.RemoveObject(ctx, bucketName, objectKey, minioOpts); err != nil {
		return fmt.Errorf(
			"removing object from storage: %w",
			apperrors.NewErrInternal().WithCause(err),
		)
	}
	return nil
}
