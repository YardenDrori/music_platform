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

func toMinioCompletePart(p ...CompletedPart) []minio.CompletePart {
	var parts []minio.CompletePart
	for _, part := range p {
		minioPart := minio.CompletePart{
			PartNumber: part.PartNumber,
			ETag:       part.ETag,
		}
		switch part.ChecksumAlgo {
		case ChecksumSha256:
			minioPart.ChecksumSHA256 = part.ChecksumValue
		}
		parts = append(parts, minioPart)
	}
	return parts
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
		return nil, fmt.Errorf(
			"generating a presigned url for a put request: %w",
			apperrors.NewErrInternal().WithCause(err),
		)
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
		return "", fmt.Errorf(
			"initiating new presigned multipart upload: %w",
			apperrors.NewErrInternal().WithCause(err),
		)
	}
	return uploadID, nil
}

func (s *service) GetPresignedMultipartPartsURLs(
	ctx context.Context,
	bucketName string,
	objectKey string,
	uploadID string,
	totalPartsCount int,
	checksums ...string,
) ([]string, error) {

	if len(checksums) != 0 && len(checksums) != totalPartsCount {
		panic(
			fmt.Sprintf(
				"called GetPresignedMultipartPartsURLs with totalPartsCount: %d, but %d items in checksums",
				totalPartsCount,
				len(checksums),
			),
		)
	}

	var presignedURLs []string
	for idx := range totalPartsCount {
		req := url.Values{
			"partNumber": {strconv.Itoa(idx + 1)},
			"uploadId":   {uploadID},
		}

		headers := http.Header{}

		if len(checksums) != 0 {
			// checksums is a slice of base64 strings
			headers.Add("x-amz-checksum-sha256", checksums[idx])
		}

		presignedURL, err := s.s3Core.PresignHeader(
			ctx,
			http.MethodPut,
			bucketName,
			objectKey,
			s.preginedURLsDuration,
			req,
			headers,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"generating presigned multipart upload URLs: %w",
				apperrors.NewErrInternal().WithCause(err),
			)
		}
		presignedURLs = append(presignedURLs, presignedURL.String())
	}

	return presignedURLs, nil
}

func (s *service) CompleteMultipartUpload(
	ctx context.Context,
	bucketName string,
	objectKey string,
	uploadID string,
	parts []CompletedPart,
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
		toMinioCompletePart(parts...),
		minioOpts,
	)
	if err != nil {
		anotherErr := s.AbortMultipartUpload(context.Background(), bucketName, objectKey, uploadID)
		if anotherErr != nil {
			return fmt.Errorf(
				"completing presgined multipart upload: %w %w",
				apperrors.NewErrInternal().WithCause(err),
				apperrors.NewErrInternal().WithCause(anotherErr),
			)
		}
		return fmt.Errorf(
			"completing presgined multipart upload: %w",
			apperrors.NewErrInternal().WithCause(err),
		)
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
		return fmt.Errorf("aborting multipart upload %w", apperrors.NewErrInternal().WithCause(err))
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
