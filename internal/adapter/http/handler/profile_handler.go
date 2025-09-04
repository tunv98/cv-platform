package handler

import (
	"cv-platform/internal/adapter/http/middleware"
	"cv-platform/internal/adapter/response"
	"cv-platform/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	uc *usecase.ProfileStoreUC
}

func NewProfileHandler(uc *usecase.ProfileStoreUC) *ProfileHandler {
	return &ProfileHandler{uc: uc}
}

type profileReq struct {
	Phone  string `uri:"id" binding:"required"`
	Status bool   `query:"status" `
}

type profileResp struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	// Option 1: Use simple logger (recommended for simple cases)
	log := middleware.SimpleLoggerFromContext(c)

	log.Info("getting profile request")

	var req profileReq
	if err := c.ShouldBindUri(&req); err != nil {
		log.Warnf("validation failed: %v", err)
		response.RespondValidationErr(c, err.Error())
		return
	}

	log.Infof("processing get profile request for phone: %s", req.Phone)

	res, err := h.uc.GetProfile(c.Request.Context(), usecase.GetProfileCmd{
		Phone: req.Phone,
	})
	if err != nil {
		log.Errorf("failed to get profile for phone %s: %v", req.Phone, err)
		response.RespondInternalErr(c, err.Error())
		return
	}

	resp := profileResp{
		ID:        res.ID,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Email:     res.Email,
		Phone:     res.Phone,
	}

	log.Infof("profile retrieved successfully: id=%s, phone=%s", res.ID, res.Phone)

	response.RespondSuccess(c, http.StatusOK, resp)
}
