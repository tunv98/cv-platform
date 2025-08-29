package port

import (
	"time"
)

type SignedURLOptions struct {
	Method      string
	ContentType string
	ExpiredAt   time.Time
}

type BlobStorage interface {
	SignedURL(objectPath string, opts SignedURLOptions) (string, error)
	Head(objectPath string) (exists bool, size int64, contentType string, err error)
}
