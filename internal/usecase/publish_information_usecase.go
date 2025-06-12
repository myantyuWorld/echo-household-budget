package usecase

import (
	domainmodel "echo-household-budget/internal/domain/model"
	"echo-household-budget/internal/domain/repository"
	domainservice "echo-household-budget/internal/domain/service"

	"github.com/davecgh/go-spew/spew"
)

type (
	PublishInformationInput struct {
		ID int
	}

	PublishInformationOutput struct {
		ID int
	}

	PublishInformationUsecase interface {
		Execute(input PublishInformationInput) (PublishInformationOutput, error)
	}

	publishInformationUsecase struct {
		informationRepository     repository.InformationRepository
		userInformationRepository repository.UserInformationRepository
		userAccountService        domainservice.UserAccountService
	}
)

// Execute implements PublishInformationUsecase.
func (p *publishInformationUsecase) Execute(input PublishInformationInput) (PublishInformationOutput, error) {
	userAccounts, err := p.userAccountService.FetchAllUserAccount()
	if err != nil {
		return PublishInformationOutput{}, err
	}

	var userInformations []*domainmodel.UserInformation
	for _, userAccount := range userAccounts {
		userInformations = append(userInformations, domainmodel.NewUserInformation(input.ID, userAccount.ID))
	}

	spew.Dump(input.ID)

	err = p.informationRepository.UpdateStatusPublished(input.ID)
	if err != nil {
		return PublishInformationOutput{}, err
	}

	for _, userInformation := range userInformations {
		err := p.userInformationRepository.Create(userInformation)
		if err != nil {
			return PublishInformationOutput{}, err
		}
	}

	return PublishInformationOutput{
		ID: input.ID,
	}, nil
}

func NewPublishInformationUsecase(informationRepository repository.InformationRepository, userInformationRepository repository.UserInformationRepository, userAccountService domainservice.UserAccountService) PublishInformationUsecase {
	return &publishInformationUsecase{
		informationRepository:     informationRepository,
		userInformationRepository: userInformationRepository,
		userAccountService:        userAccountService,
	}
}
