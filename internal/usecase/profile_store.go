package usecase

import (
	"context"
	logger "cv-platform/internal/log"
	"fmt"
)

type ProfileStoreUC struct {
}

func NewProfileStoreUC() *ProfileStoreUC {
	return &ProfileStoreUC{}
}

type GetProfileCmd struct {
	Phone string
}

type GetProfileResult struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

func (uc *ProfileStoreUC) GetProfile(ctx context.Context, cmd GetProfileCmd) (*GetProfileResult, error) {
	// Option 2: Use simple logger for usecase
	log := logger.SimpleFromContext(ctx)
	log.Infof("getting profile for phone: %s", cmd.Phone)

	if cmd.Phone == "1111" {
		log.Warnf("profile not found for phone %s: blacklisted", cmd.Phone)
		return nil, fmt.Errorf("phone number %s is not found", cmd.Phone)
	}

	result := &GetProfileResult{
		ID:        "1",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "John@Doe.gmail.com",
		Phone:     cmd.Phone,
	}

	log.Infof("profile retrieved successfully: id=%s, phone=%s, email=%s",
		result.ID, result.Phone, result.Email)

	return result, nil
}
