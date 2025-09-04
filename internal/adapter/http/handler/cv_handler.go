package handler

import (
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
	var req startReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.RespondValidationErr(c, err.Error())
		return
	}

	res, err := h.uc.StartUpload(usecase.StartUploadCmd{
		FileName: req.FileName,
		MimeType: req.MimeType,
	})
	if err != nil {
		response.RespondInternalErr(c, err.Error())
		return
	}

	resp := startResp{
		ID:        res.ID,
		ObjectKey: res.ObjectKey,
		SignedURL: res.SignedURL,
		ExpiredAt: res.ExpiredAt,
	}

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
	id := c.Param("id")
	cv, err := h.uc.CompleteUpload(usecase.CompleteUploadCmd{ID: id})
	if err != nil {
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

	response.RespondSuccess(c, http.StatusOK, resp)
}
