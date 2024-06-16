package images

import "mime/multipart"

type uploadFile struct {
	UserID string
	File   *multipart.FileHeader
}
