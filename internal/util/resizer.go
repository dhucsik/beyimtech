package util

import (
	"image"

	"github.com/disintegration/imaging"
)

func Resize(source image.Image, width int) image.Image {
	return imaging.Resize(source, width, 0, imaging.Lanczos)
}
