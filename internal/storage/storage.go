package storage

import (
	"beyimtech-test/configs"
	awss3 "beyimtech-test/internal/storage/aws-s3"
	"context"
	"io"
)

type Storage interface {
	UploadFile(ctx context.Context, body io.ReadSeeker, filename string) (string, error)
	DeleteFile(ctx context.Context, filename string) error
}

func NewStorage(cfg *configs.Config) (Storage, error) {
	stg, err := awss3.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return stg, nil
}
