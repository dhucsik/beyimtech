package images

import (
	"beyimtech-test/internal/entity"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) SaveImages(ctx context.Context, images []*entity.Image) error {
	if len(images) == 0 {
		return nil
	}

	mm := convertArray(images)
	batch := &pgx.Batch{}

	for _, m := range mm {
		batch.Queue(insertStmt, m.EntityID, m.EntityType, m.ImageURL, m.ContentType, m.Filename, m.Format)
	}

	results := r.db.SendBatch(ctx, batch)
	return results.Close()
}

func (r *Repository) GetImagesByUserID(ctx context.Context, userID string) ([]*entity.Image, error) {
	rows, err := r.db.Query(ctx, selectUserStmt, userID)
	if err != nil {
		return nil, err
	}

	var mm models
	for rows.Next() {
		m := &model{}

		if err := rows.Scan(&m.EntityID, &m.EntityType, &m.ImageURL, &m.ContentType, &m.Filename, &m.Format); err != nil {
			return nil, err
		}

		mm = append(mm, m)
	}

	return mm.convertImages(), nil
}

func (r *Repository) GetImages(ctx context.Context, limit, offset int) ([]*entity.Image, error) {
	rows, err := r.db.Query(ctx, selectStmt, limit, offset)
	if err != nil {
		return nil, err
	}

	var mm models
	for rows.Next() {
		m := &model{}

		if err := rows.Scan(&m.EntityID, &m.EntityType, &m.ImageURL, &m.ContentType, &m.Filename, &m.Format); err != nil {
			return nil, err
		}

		mm = append(mm, m)
	}

	return mm.convertImages(), nil
}

func (r *Repository) GetByURL(ctx context.Context, url string) (*entity.Image, error) {
	m := &model{}

	row := r.db.QueryRow(ctx, getByURLStmt, url)
	if err := row.Scan(&m.EntityID, &m.EntityType, &m.ImageURL, &m.ContentType, &m.Filename, &m.Format); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return m.convertToImage(), nil
}

func (r *Repository) DeleteByURL(ctx context.Context, url string) error {
	_, err := r.db.Exec(ctx, deleteByURLStmt, url)
	return err
}

func (r *Repository) DeleteByUserID(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx, deleteByUserIDStmt, userID)
	return err
}
