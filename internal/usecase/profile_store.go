package usecase

import (
	"context"
	"cv-platform/pkg/logger"
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
	log := logger.FLogFromContext(ctx)
	log.Infof("getting profile for phone: %s", cmd.Phone)

	if cmd.Phone == "****" {
		panic("simulated panic for phone ****")
	}
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

type CreateProfileCmd struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Age       int
	Gender    string
}

type CreateProfileResult struct {
	ID string
}

func (uc *ProfileStoreUC) CreateProfile(ctx context.Context, cmd CreateProfileCmd) (*CreateProfileResult, error) {
	log := logger.FLogFromContext(ctx)
	log.Infof("creating profile: first_name=%s, last_name=%s, email=%s, phone=%s, age=%d, gender=%s",
		cmd.FirstName, cmd.LastName, cmd.Email, cmd.Phone, cmd.Age, cmd.Gender)

	return &CreateProfileResult{ID: "1"}, nil
}
