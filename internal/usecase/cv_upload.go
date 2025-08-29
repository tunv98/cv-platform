package usecase

import (
	"cv-platform/internal/domain"
	logger "cv-platform/internal/log"
	"cv-platform/internal/port"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CVUploadUC struct {
	storage port.BlobStorage
	repo    port.CVRepository
}

func NewCVUploadUC(storage port.BlobStorage, repo port.CVRepository) *CVUploadUC {
	return &CVUploadUC{
		storage: storage,
		repo:    repo,
	}
}

type StartUploadCmd struct {
	FileName string
	MimeType string
}

type StartUploadResult struct {
	ID        string
	ObjectKey string
	SignedURL string
	ExpiredAt time.Time
}

func (uc *CVUploadUC) StartUpload(cmd StartUploadCmd) (*StartUploadResult, error) {
	id := uuid.New().String()
	var ext string
	if dot := lastDot(cmd.FileName); dot != -1 {
		ext = cmd.FileName[dot+1:]
	}

	objectKey := fmt.Sprintf("cv/%s.%s", id, ext)
	logger.L().Info("generating signed url", zap.String("id", id), zap.String("object_key", objectKey), zap.String("mime_type", cmd.MimeType))

	opts := port.SignedURLOptions{
		Method:      "PUT",
		ContentType: cmd.MimeType,
		ExpiredAt:   time.Now().Add(10 * time.Minute),
	}

	url, err := uc.storage.SignedURL(objectKey, opts)
	if err != nil {
		logger.L().Error("failed to get signed url", zap.Error(err))
		return nil, err
	}

	cv := &domain.CV{
		ID:        id,
		FileName:  cmd.FileName,
		MimeType:  cmd.MimeType,
		Size:      0,
		GCSPath:   objectKey,
		Status:    domain.CVStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.repo.Create(cv); err != nil {
		logger.L().Error("failed to create cv", zap.Error(err))
		return nil, err
	}

	res := &StartUploadResult{
		ID:        id,
		ObjectKey: objectKey,
		SignedURL: url,
		ExpiredAt: opts.ExpiredAt,
	}

	logger.L().Info("upload initialized", zap.String("id", id), zap.String("object_key", objectKey), zap.Time("expires_at", res.ExpiredAt))
	return res, nil
}

type CompleteUploadCmd struct {
	ID string
}

func (uc *CVUploadUC) CompleteUpload(cmd CompleteUploadCmd) (*domain.CV, error) {
	logger.L().Info("completing upload", zap.String("id", cmd.ID))
	cv, err := uc.repo.FindByID(cmd.ID)
	if err != nil {
		logger.L().Error("failed to find cv", zap.Error(err))
		return nil, err
	}
	ok, size, ctype, err := uc.storage.Head(cv.GCSPath)
	if err != nil {
		logger.L().Error("failed to head cv", zap.Error(err))
		return nil, err
	}
	if !ok {
		logger.L().Error("object not found", zap.String("object_key", cv.GCSPath))
		return nil, fmt.Errorf("object not found: %s", cv.GCSPath)
	}
	cv.Size = size
	if ctype != "" {
		cv.MimeType = ctype
	}
	cv.Status = domain.CVStatusUploaded
	cv.UpdatedAt = time.Now()
	if err := uc.repo.Update(cv); err != nil {
		logger.L().Error("failed to update cv", zap.Error(err))
		return nil, fmt.Errorf("failed to update cv: %w", err)
	}
	logger.L().Info("upload completed", zap.String("id", cmd.ID))
	return cv, nil
}

func lastDot(s string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			return i
		}
	}
	return -1
}
