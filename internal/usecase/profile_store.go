package usecase

import "fmt"

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

func (uc *ProfileStoreUC) GetProfile(cmd GetProfileCmd) (*GetProfileResult, error) {
	if cmd.Phone == "1111" {
		return nil, fmt.Errorf("phone number %s is not found", cmd.Phone)
	}
	return &GetProfileResult{
		ID:        "1",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "John@Doe.gmail.com",
		Phone:     cmd.Phone,
	}, nil
}
