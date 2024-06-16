package images

import (
	"beyimtech-test/internal/entity"

	"github.com/samber/lo"
)

type model struct {
	EntityID    string
	EntityType  string
	ImageURL    string
	ContentType string
	Filename    string
	Format      int
}

func (m *model) convertToImage() *entity.Image {
	return &entity.Image{
		ID:         m.EntityID,
		EntityType: m.EntityType,
		URL:        m.ImageURL,
		Filename:   m.Filename,
		Format:     m.Format,
		Filetype:   m.ContentType,
	}
}

type models []*model

func (mm models) convertImages() []*entity.Image {
	return lo.Map(mm, func(m *model, _ int) *entity.Image {
		return m.convertToImage()
	})
}

func convertArray(items []*entity.Image) models {
	return lo.Map(items, func(item *entity.Image, _ int) *model {
		return &model{
			EntityID:    item.ID,
			EntityType:  item.EntityType,
			ImageURL:    item.URL,
			ContentType: item.Filetype,
			Filename:    item.Filename,
			Format:      item.Format,
		}
	})
}
