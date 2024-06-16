package awss3

import (
	"beyimtech-test/configs"
	"context"
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Storage struct {
	s3Client    *s3.S3
	bucket      string
	endpoint    string
	cdnEndpoint string
}

func NewClient(
	cfg *configs.Config,
) (*Storage, error) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String(cfg.Env.Get("s3_endpoint")),
		Region:   aws.String(cfg.Env.Get("s3_region")),
		Credentials: credentials.NewStaticCredentials(
			cfg.Env.Get("s3_access_key"),
			cfg.Env.Get("s3_secret_key"),
			"",
		),
	})
	if err != nil {
		return nil, errors.Wrap(err, "storage.NewClient")
	}

	return &Storage{
		s3Client:    s3.New(sess),
		bucket:      cfg.Env.Get("s3_bucket"),
		endpoint:    cfg.Env.Get("s3_endpoint"),
		cdnEndpoint: cfg.Env.Get("s3_cdn_endpoint"),
	}, nil
}

func (s *Storage) UploadFile(ctx context.Context, body io.ReadSeeker, filename string) (string, error) {
	object := s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename),
		Body:   body,
		ACL:    aws.String("public-read"),
	}

	_, err := s.s3Client.PutObjectWithContext(ctx, &object)
	if err != nil {
		return "", errors.Wrapf(err, "storage.UploadFile - filename: %s", filename)
	}

	return s.genURL(filename), err
}

func (s *Storage) DeleteFile(ctx context.Context, filename string) error {
	object := s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename),
	}

	_, err := s.s3Client.DeleteObjectWithContext(ctx, &object)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) genURL(filename string) string {
	return fmt.Sprintf("https://%s.%s/%s", s.bucket, s.cdnEndpoint, filename)
}
