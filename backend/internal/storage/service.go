package storage

import (
	"context"
	"io"
	"net/url"
)

type Service interface {
	//errors:
	//fmt
	//Initiates new multipart upload and returns the new uploadID.
	InitiateMultipartUpload(
		ctx context.Context,
		bucketName string,
		objectKey string,
		opts PutOptions,
	) (string, error)

	//errors:
	//fmt
	PresignMultipartUploadPutURLs(
		ctx context.Context,
		bucketName string,
		objectKey string,
		uploadID string,
		totalPartsCount int,
		checksums ...string,
	) ([]string, error)

	//errors:
	//fmt
	CompleteMultipartUpload(
		ctx context.Context,
		bucketName string,
		objectKey string,
		uploadID string,
		parts []CompletedPart,
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

	DeleteObject(
		ctx context.Context,
		bucketName string,
		objectKey string,
		opts DeleteOptions,
	) error
}
