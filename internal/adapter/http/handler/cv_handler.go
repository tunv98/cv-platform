package handler

import (
	"cv-platform/internal/adapter/http/middleware"
	"cv-platform/internal/adapter/response"
	"cv-platform/internal/usecase"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CVHandler struct {
	uc *usecase.CVUploadUC
}

func NewCVHandler(uc *usecase.CVUploadUC) *CVHandler {
	return &CVHandler{uc: uc}
}

type startReq struct {
	FileName string `json:"file_name" biding:"required"`
	MimeType string `json:"mime_type" biding:"required"`
}

type startResp struct {
	ID        string    `json:"id"`
	ObjectKey string    `json:"object_key"`
	SignedURL string    `json:"signed_url"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (h *CVHandler) StartUpload(c *gin.Context) {
	log := middleware.SimpleLoggerFromContext(c)

	log.Info("starting upload request")

	var req startReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf("validation failed: %v", err)
		response.RespondValidationErr(c, err.Error())
		return
	}

	log.Infof("processing upload request: file=%s, type=%s", req.FileName, req.MimeType)

	res, err := h.uc.StartUpload(c.Request.Context(), usecase.StartUploadCmd{
		FileName: req.FileName,
		MimeType: req.MimeType,
	})
	if err != nil {
		log.Errorf("failed to start upload for file %s: %v", req.FileName, err)
		response.RespondInternalErr(c, err.Error())
		return
	}

	resp := startResp{
		ID:        res.ID,
		ObjectKey: res.ObjectKey,
		SignedURL: res.SignedURL,
		ExpiredAt: res.ExpiredAt,
	}

	log.Infof("upload started successfully: id=%s, expires_at=%v", res.ID, res.ExpiredAt)

	response.RespondSuccess(c, http.StatusOK, resp)
}

type completeResp struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
	GCSPath  string `json:"gcs_path"`
}

func (h *CVHandler) CompleteUpload(c *gin.Context) {
	log := middleware.SimpleLoggerFromContext(c)
	id := c.Param("id")

	log.Infof("completing upload request for id: %s", id)

	cv, err := h.uc.CompleteUpload(c.Request.Context(), usecase.CompleteUploadCmd{ID: id})
	if err != nil {
		log.Errorf("failed to complete upload for id %s: %v", id, err)
		response.RespondBadRequest(c, err.Error())
		return
	}

	resp := completeResp{
		ID:       cv.ID,
		Status:   string(cv.Status),
		Size:     cv.Size,
		MimeType: cv.MimeType,
		GCSPath:  cv.GCSPath,
	}

	log.Infof("upload completed successfully: id=%s, status=%s, size=%d", cv.ID, cv.Status, cv.Size)

	response.RespondSuccess(c, http.StatusOK, resp)
}
