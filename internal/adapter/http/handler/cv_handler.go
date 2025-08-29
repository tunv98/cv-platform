package handler

import (
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.uc.StartUpload(usecase.StartUploadCmd{
		FileName: req.FileName,
		MimeType: req.MimeType,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, startResp{
		ID: res.ID,
	})
}

func (h *CVHandler) CompleteUpload(c *gin.Context) {
	id := c.Param("id")
	cv, err := h.uc.CompleteUpload(usecase.CompleteUploadCmd{ID: id})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": cv.ID, "status": cv.Status, "size": cv.Size, "mime": cv.MimeType, "gcs_path": cv.GCSPath,
	})
}
