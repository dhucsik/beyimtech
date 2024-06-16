package repositories

import (
	"beyimtech-test/internal/entity"
	"beyimtech-test/internal/repositories/postgres"
	"context"
)

type ImageRepository interface {
	SaveImages(ctx context.Context, images []*entity.Image) error
	GetImagesByUserID(ctx context.Context, userID string) ([]*entity.Image, error)
	GetImages(ctx context.Context, limit, offset int) ([]*entity.Image, error)
	DeleteByUserID(ctx context.Context, userID string) error
	DeleteByURL(ctx context.Context, url string) error
	GetByURL(ctx context.Context, url string) (*entity.Image, error)
}

func NewRepository(ctx context.Context, dsn string) (ImageRepository, error) {
	return initPostgre(ctx, dsn)
}

func initPostgre(ctx context.Context, dsn string) (ImageRepository, error) {
	db, err := postgres.Dial(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return postgres.NewImagesRepo(db), nil
}
