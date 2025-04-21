package domainservice

import (
	domainmodel "echo-household-budget/internal/domain/model"
)

type HouseHoldService interface {
	FetchHouseHold(houseHoldID domainmodel.HouseHoldID) (*domainmodel.HouseHold, error)
	ShareHouseHold(houseHoldID domainmodel.HouseHoldID, inviteUserID domainmodel.UserID) error
}

type houseHoldService struct {
	houseHoldRepository domainmodel.HouseHoldRepository
}

// FetchHouseHold implements HouseHoldService.
func (h *houseHoldService) FetchHouseHold(houseHoldID domainmodel.HouseHoldID) (*domainmodel.HouseHold, error) {
	houseHold, err := h.houseHoldRepository.FindByHouseHoldID(houseHoldID)
	if err != nil {
		return nil, err
	}
	return houseHold, nil
}

// ShareHouseHold implements HouseHoldService.
func (h *houseHoldService) ShareHouseHold(houseHoldID domainmodel.HouseHoldID, inviteUserID domainmodel.UserID) error {
	userHouseHold := &domainmodel.UserHouseHold{
		HouseHoldID: houseHoldID,
		UserID:      inviteUserID,
	}

	return h.houseHoldRepository.CreateUserHouseHold(userHouseHold)
}

func NewHouseHoldService(houseHoldRepository domainmodel.HouseHoldRepository) HouseHoldService {
	return &houseHoldService{
		houseHoldRepository: houseHoldRepository,
	}
}
