package storage

import (
	"context"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/s3utils"

	"github.com/YardenDrori/music-platform/internal/apperrors"
)

type service struct {
	s3Core *minio.Core
}

func NewService(client *minio.Client) Service {
	return &service{s3Core: &minio.Core{Client: client}}
}

func (s *service) InitiateMutlipartPresignedUpload(
	ctx context.Context,
	bucketName string,
	objectKey string,
) (string error) {
	s.s3Core.NewMultipartUpload()

	panic("")
}
