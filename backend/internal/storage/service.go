package storage

import (
	"context"
	"io"
	"net/url"

	"github.com/minio/minio-go/v7"
)

type PutOptions struct {
	ContentType    string
	SendContentMD5 bool
}

type Service interface {
	//errors:
	//fmt
	InitiateMultipartUpload(
		ctx context.Context,
		bucketName string,
		objectKey string,
		opts PutOptions,
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
		opts PutOptions,
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
	) (*url.URL, error)

	PutObject(
		ctx context.Context,
		bucketName string,
		objectKey string,
		reader io.Reader,
		size int64,
		opts PutOptions,
	) error

	BuildPublicGetUrl(bucketName string, objectKey string) *string
}
