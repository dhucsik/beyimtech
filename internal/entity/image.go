package entity

import (
	"bytes"
	"image"
)

type Image struct {
	ID         string        `json:"id"`
	EntityType string        `json:"entity_type"`
	Src        *bytes.Reader `json:"-"`
	URL        string        `json:"url"`
	Filename   string        `json:"filename"`
	Format     int           `json:"format"`
	Filetype   string        `json:"filetype"`
}

type ImageParams struct {
	FileType string
	UserID   string
	Image    image.Image
	Filename string
	Format   int
}
