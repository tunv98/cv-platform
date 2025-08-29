package domain

import "time"

type CVStatus string

const (
	CVStatusPending  CVStatus = "pending"
	CVStatusUploaded CVStatus = "uploaded"
)

type CV struct {
	ID        string
	FileName  string
	MimeType   string
	Size      int64
	GCSPath   string
	Status    CVStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}
