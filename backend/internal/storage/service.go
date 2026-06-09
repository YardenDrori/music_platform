package storage

import (
	"context"
	"net/url"

	"github.com/minio/minio-go/v7"
)

type Service interface {
	//errors:
	//fmt
	InitiateMultipartUpload(
		ctx context.Context,
		bucketName string,
		objectKey string,
		opts minio.PutObjectOptions,
	) (string, error)

	//errors:
	//fmt
	GetPresignedMultipartPartsURLs(
		ctx context.Context,
		bucketName string,
		objectKey string,
		uploadID string,
		totalPartsCount int,
	) ([]string, error)

	//errors:
	//fmt
	CompleteMultipartUpload(
		ctx context.Context,
		bucketName string,
		objectKey string,
		uploadID string,
		ETags []minio.CompletePart,
		opts minio.PutObjectOptions,
	) error

	//errors:
	//fmt
	AbortMultipartUpload(
		ctx context.Context,
		bucketName string,
		objectKey string,
		uploadID string,
	) error

	//erros:
	//fmt
	PresignedUpload(
		ctx context.Context,
		bucketName string,
		objectKey string,
		opts minio.PutObjectOptions,
	) (*url.URL, error)
}
