package handler

import (
	"cv-platform/internal/adapter/http/middleware"
	"cv-platform/internal/adapter/response"
	"cv-platform/internal/metrics"
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

type getProfileReq struct {
	Phone  string `uri:"id" binding:"required"`
	Status bool   `query:"status" `
}

type getProfileResp struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	// Option 1: Use simple logger (recommended for simple cases)
	log := middleware.LoggerFromContext(c)

	log.Info("getting profile request")

	var req getProfileReq
	if err := c.ShouldBindUri(&req); err != nil {
		log.Warnf("validation failed: %v", err)
		metrics.RecordProfileRequest("get", "validation_error")
		response.RespondValidationErr(c, err.Error())
		return
	}

	log.Infof("processing get profile request for phone: %s", req.Phone)

	res, err := h.uc.GetProfile(c.Request.Context(), usecase.GetProfileCmd{
		Phone: req.Phone,
	})
	if err != nil {
		log.Errorf("failed to get profile for phone %s: %v", req.Phone, err)
		metrics.RecordProfileRequest("get", "error")
		response.RespondInternalErr(c, err.Error())
		return
	}

	resp := getProfileResp{
		ID:        res.ID,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Email:     res.Email,
		Phone:     res.Phone,
	}

	log.Infof("profile retrieved successfully: id=%s, phone=%s", res.ID, res.Phone)
	metrics.RecordProfileRequest("get", "success")

	response.RespondSuccess(c, http.StatusOK, resp)
}

type createProfileReq struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone" binding:"required"`
	Age       int    `json:"age" binding:"gte=15,lte=130"`
	Gender    string `json:"gender" binding:"oneof=male female"`
}

type createProfileResp struct {
	ID string `json:"id"`
}

func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	log := middleware.LoggerFromContext(c)

	log.Info("creating profile request")

	var req createProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf("validation failed: %v", err)
		metrics.RecordProfileRequest("create", "validation_error")
		response.RespondValidationErr(c, err.Error())
		return
	}

	log.Infof("processing create profile request: first_name=%s, last_name=%s, email=%s, phone=%s, age=%d, gender=%s",
		req.FirstName, req.LastName, req.Email, req.Phone, req.Age, req.Gender)

	res, err := h.uc.CreateProfile(c.Request.Context(), usecase.CreateProfileCmd{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Age:       req.Age,
		Gender:    req.Gender,
	})

	if err != nil {
		log.Errorf("failed to create profile: %v", err)
		metrics.RecordProfileRequest("create", "error")
		response.RespondInternalErr(c, err.Error())
		return
	}

	log.Infof("profile created successfully: id=%s", res.ID)
	metrics.RecordProfileRequest("create", "success")

	response.RespondSuccess(c, http.StatusOK, createProfileResp{ID: res.ID})
}
