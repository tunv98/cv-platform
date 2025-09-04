package usecase

import (
	"context"
	"cv-platform/internal/domain"
	logger "cv-platform/internal/log"
	"cv-platform/internal/port"
	"fmt"
	"time"

	"github.com/google/uuid"
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

func (uc *CVUploadUC) StartUpload(ctx context.Context, cmd StartUploadCmd) (*StartUploadResult, error) {
	log := logger.SimpleFromContext(ctx)
	log.Infof("starting upload process: file=%s, type=%s", cmd.FileName, cmd.MimeType)

	id := uuid.New().String()
	var ext string
	if dot := lastDot(cmd.FileName); dot != -1 {
		ext = cmd.FileName[dot+1:]
	}

	objectKey := fmt.Sprintf("cv/%s.%s", id, ext)
	log.Infof("generating signed url: id=%s, key=%s, ext=%s", id, objectKey, ext)

	opts := port.SignedURLOptions{
		Method:      "PUT",
		ContentType: cmd.MimeType,
		ExpiredAt:   time.Now().Add(10 * time.Minute),
	}

	url, err := uc.storage.SignedURL(objectKey, opts)
	if err != nil {
		log.Errorf("failed to get signed url for id %s: %v", id, err)
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

	log.Infof("saving cv to repository: id=%s, status=%s", id, cv.Status)

	if err := uc.repo.Create(cv); err != nil {
		log.Errorf("failed to create cv for id %s: %v", id, err)
		return nil, err
	}

	res := &StartUploadResult{
		ID:        id,
		ObjectKey: objectKey,
		SignedURL: url,
		ExpiredAt: opts.ExpiredAt,
	}

	log.Infof("upload initialized successfully: id=%s, expires_at=%v", id, res.ExpiredAt)
	return res, nil
}

type CompleteUploadCmd struct {
	ID string
}

func (uc *CVUploadUC) CompleteUpload(ctx context.Context, cmd CompleteUploadCmd) (*domain.CV, error) {
	log := logger.SimpleFromContext(ctx)
	log.Infof("completing upload process for id: %s", cmd.ID)

	cv, err := uc.repo.FindByID(cmd.ID)
	if err != nil {
		log.Errorf("failed to find cv for id %s: %v", cmd.ID, err)
		return nil, err
	}

	log.Infof("checking object in storage: path=%s", cv.GCSPath)

	ok, size, ctype, err := uc.storage.Head(cv.GCSPath)
	if err != nil {
		log.Errorf("failed to head cv for id %s at path %s: %v", cmd.ID, cv.GCSPath, err)
		return nil, err
	}
	if !ok {
		log.Errorf("object not found in storage: id=%s, path=%s", cmd.ID, cv.GCSPath)
		return nil, fmt.Errorf("object not found: %s", cv.GCSPath)
	}

	log.Infof("updating cv with file information: id=%s, size=%d, type=%s", cmd.ID, size, ctype)

	cv.Size = size
	if ctype != "" {
		cv.MimeType = ctype
	}
	cv.Status = domain.CVStatusUploaded
	cv.UpdatedAt = time.Now()

	if err := uc.repo.Update(cv); err != nil {
		log.Errorf("failed to update cv for id %s: %v", cmd.ID, err)
		return nil, fmt.Errorf("failed to update cv: %w", err)
	}

	log.Infof("upload completed successfully: id=%s, status=%s, size=%d", cmd.ID, cv.Status, cv.Size)
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
