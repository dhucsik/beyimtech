package images

import (
	"beyimtech-test/internal/enum"
	apiErrors "beyimtech-test/internal/errors"
)

const maxFileSize = 5 * 1024 * 1024

func (r *uploadImageRequest) Validate() error {
	if r.file.Size > maxFileSize {
		return apiErrors.ErrFileSizeExceeds
	}

	contentType := r.file.Header.Get("Content-Type")

	if contentType != enum.MimeTypeImageJPEG && contentType != enum.MimeTypeImagePng &&
		contentType != enum.MimeTypeImageWebp {
		return apiErrors.ErrNotSupportedImgFormat
	}

	return nil
}
