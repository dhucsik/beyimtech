package converter

import (
	"bytes"
	"context"
	"image"

	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
)

type Converter interface {
	ConvertImgToWebp(ctx context.Context, img image.Image) (*bytes.Reader, error)
}

type converter struct {
	webpQuality int
}

func NewConverter(webpQuality int) Converter {
	return &converter{
		webpQuality: webpQuality,
	}
}

func (c *converter) ConvertImgToWebp(_ context.Context, img image.Image) (*bytes.Reader, error) {
	out := new(bytes.Buffer)

	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, float32(c.webpQuality))
	if err != nil {
		return nil, err
	}

	err = webp.Encode(out, img, options)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(out.Bytes()), nil
}
