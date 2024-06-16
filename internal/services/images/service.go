package images

import (
	"beyimtech-test/internal/entity"
	"beyimtech-test/internal/enum"
	apiErrors "beyimtech-test/internal/errors"
	"beyimtech-test/internal/repositories"
	"beyimtech-test/internal/storage"
	"beyimtech-test/internal/util"
	"beyimtech-test/internal/util/converter"
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"sync"

	"github.com/kolesa-team/go-webp/decoder"
	"github.com/kolesa-team/go-webp/webp"
	"golang.org/x/sync/errgroup"
)

type Service interface {
	SendFile(ctx context.Context, userID string, file *multipart.FileHeader)
	ReceiveFiles(ctx context.Context)
	GetImagesByUserID(ctx context.Context, userID string) ([]*entity.Image, error)
	GetImages(ctx context.Context, limit, offset int) ([]*entity.Image, error)
	DeleteImagesByUserID(ctx context.Context, userID string) error
	DeleteImageByURL(ctx context.Context, imageURL string) error
}

type service struct {
	storage    storage.Storage
	converter  converter.Converter
	repository repositories.ImageRepository
	filesChan  chan uploadFile
	mtx        sync.Mutex
}

func NewService(
	storage storage.Storage,
	converter converter.Converter,
	repository repositories.ImageRepository,
) Service {
	return &service{
		storage:    storage,
		converter:  converter,
		repository: repository,
		filesChan:  make(chan uploadFile, 100),
		mtx:        sync.Mutex{},
	}
}

func (s *service) GetImagesByUserID(ctx context.Context, userID string) ([]*entity.Image, error) {
	return s.repository.GetImagesByUserID(ctx, userID)
}

func (s *service) GetImages(ctx context.Context, limit, offset int) ([]*entity.Image, error) {
	return s.repository.GetImages(ctx, limit, offset)
}

func (s *service) DeleteImagesByUserID(ctx context.Context, userID string) error {
	images, err := s.repository.GetImagesByUserID(ctx, userID)
	if err != nil {
		return err
	}

	err = s.repository.DeleteByUserID(ctx, userID)
	if err != nil {
		return err
	}

	for _, img := range images {
		err := s.storage.DeleteFile(ctx, img.Filename)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) DeleteImageByURL(ctx context.Context, imageURL string) error {
	image, err := s.repository.GetByURL(ctx, imageURL)
	if err != nil {
		return err
	}

	if image == nil {
		return nil
	}

	err = s.storage.DeleteFile(ctx, image.Filename)
	if err != nil {
		return err
	}

	return s.repository.DeleteByURL(ctx, imageURL)
}

func (s *service) SendFile(_ context.Context, userID string, file *multipart.FileHeader) {
	s.filesChan <- uploadFile{
		File:   file,
		UserID: userID,
	}
}

func (s *service) ReceiveFiles(ctx context.Context) {
	for file := range s.filesChan {
		err := s.UploadUserImage(ctx, file.UserID, file.File)
		if err != nil {
			log.Println(err)
		}
	}
}

func (s *service) ResizeAndUploadImage(ctx context.Context, params *entity.ImageParams) (*entity.Image, error) {
	img := params.Image

	if params.Format > 0 {
		img = util.Resize(params.Image, params.Format)
	}

	webpImg, err := s.converter.ConvertImgToWebp(ctx, img)
	if err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("%s/%s_%d.webp", params.UserID, params.Filename, params.Format)
	url, err := s.UploadImage(ctx, filename, webpImg)
	if err != nil {
		return nil, err
	}

	return &entity.Image{
		ID:         params.UserID,
		EntityType: "user",
		URL:        url,
		Src:        webpImg,
		Filename:   filename,
		Format:     params.Format,
		Filetype:   enum.MimeTypeImageWebp,
	}, nil
}

func (s *service) UploadUserImage(ctx context.Context, userID string, file *multipart.FileHeader) error {
	filetype := file.Header.Get("Content-Type")
	reader, err := file.Open()
	if err != nil {
		return err
	}
	defer reader.Close()

	img, err := decode(reader, filetype)
	if err != nil {
		return err
	}

	imagesToSave := make([]*entity.Image, 0, 5)
	url, err := s.UploadImg(ctx, img, filetype, fmt.Sprintf("%s/%s", userID, file.Filename))
	if err != nil {
		return err
	}

	imagesToSave = append(imagesToSave, &entity.Image{
		ID:         userID,
		EntityType: "user",
		Filename:   fmt.Sprintf("%s/%s", userID, file.Filename),
		Filetype:   filetype,
		URL:        url,
	})
	ext := filepath.Ext(file.Filename)
	filename := strings.TrimSuffix(file.Filename, ext)

	g, contx := errgroup.WithContext(ctx)
	for _, format := range []int{0, 250, 500, 750} {
		format := format
		g.Go(func() error {
			uploaded, err := s.ResizeAndUploadImage(contx, &entity.ImageParams{
				UserID:   userID,
				Filename: filename,
				Image:    img,
				Format:   format,
			})

			s.mtx.Lock()
			imagesToSave = append(imagesToSave, uploaded)
			s.mtx.Unlock()

			return err
		})
	}

	err = g.Wait()
	if err != nil {
		return err
	}

	return s.repository.SaveImages(ctx, imagesToSave)
}

func (s *service) UploadImage(ctx context.Context, filename string, reader *bytes.Reader) (string, error) {
	url, err := s.storage.UploadFile(ctx, reader, filename)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (s *service) UploadImg(ctx context.Context, img image.Image, filetype, filename string) (string, error) {
	out := new(bytes.Buffer)

	switch filetype {
	case enum.MimeTypeImageJPEG:
		err := jpeg.Encode(out, img, nil)
		if err != nil {
			return "", err
		}

	case enum.MimeTypeImagePng:
		err := png.Encode(out, img)
		if err != nil {
			return "", err
		}

	default:
		return "", apiErrors.ErrNotSupportedImgFormat
	}

	return s.storage.UploadFile(ctx, bytes.NewReader(out.Bytes()), filename)
}

func decode(reader io.Reader, filetype string) (image.Image, error) {
	switch filetype {
	case enum.MimeTypeImageJPEG:
		return jpeg.Decode(reader)

	case enum.MimeTypeImagePng:
		return png.Decode(reader)

	case enum.MimeTypeImageWebp:
		return webp.Decode(reader, &decoder.Options{})

	default:
		return nil, apiErrors.ErrNotSupportedImgFormat
	}
}
