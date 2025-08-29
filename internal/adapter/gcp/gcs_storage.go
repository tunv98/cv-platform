// internal/adapter/gcp/gcs_storage.go
package gcp

import (
	"context"
	"errors"
	"net/http"
	"time"

	"cv-platform/internal/port"

	logger "cv-platform/internal/log"

	"cloud.google.com/go/storage"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

type GCSStorage struct {
	client      *storage.Client
	bucket      string
	signerEmail string
	privateKey  []byte
}

func NewGCSStorage(ctx context.Context, bucket string, credsJSON []byte) (*GCSStorage, error) {
	var (
		cl  *storage.Client
		err error
	)
	if len(credsJSON) > 0 {
		cl, err = storage.NewClient(ctx, option.WithCredentialsJSON(credsJSON))
	} else {
		cl, err = storage.NewClient(ctx)
	}
	if err != nil {
		logger.L().Error("failed to create gcs client", zap.Error(err))
		return nil, err
	}
	return &GCSStorage{client: cl, bucket: bucket}, nil
}

func (g *GCSStorage) SignedURL(object string, opts port.SignedURLOptions) (string, error) {
	return storage.SignedURL(g.bucket, object, &storage.SignedURLOptions{
		Scheme:      storage.SigningSchemeV4,
		Method:      http.MethodPut,
		Expires:     opts.ExpiredAt,
		ContentType: opts.ContentType,
	})
}

func (g *GCSStorage) Head(object string) (bool, int64, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	attrs, err := g.client.Bucket(g.bucket).Object(object).Attrs(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return false, 0, "", nil
		}
		return false, 0, "", err
	}
	return true, attrs.Size, attrs.ContentType, nil
}
