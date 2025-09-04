package handler

import (
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
	var req profileReq
	if err := c.ShouldBindUri(&req); err != nil {
		response.RespondValidationErr(c, err.Error())
		return
	}

	res, err := h.uc.GetProfile(usecase.GetProfileCmd{
		Phone: req.Phone,
	})
	if err != nil {
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
	response.RespondSuccess(c, http.StatusOK, resp)
}
