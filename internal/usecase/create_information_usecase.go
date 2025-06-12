package usecase

import (
	domainmodel "echo-household-budget/internal/domain/model"
	repository "echo-household-budget/internal/domain/repository"
)

type (
	CreateInformationInput struct {
		Title    string
		Content  string
		Category string
	}

	CreateInformationOutput struct {
		ID int
	}

	CreateInformationUsecase interface {
		Execute(input CreateInformationInput) (CreateInformationOutput, error)
	}

	createInformationUsecase struct {
		informationRepository repository.InformationRepository
	}
)

// Execute implements CreateInformationUsecase.
func (c *createInformationUsecase) Execute(input CreateInformationInput) (CreateInformationOutput, error) {
	information := domainmodel.NewInformation(input.Title, input.Content, input.Category)

	err := c.informationRepository.Create(information)
	if err != nil {
		return CreateInformationOutput{}, err
	}

	return CreateInformationOutput{
		ID: information.ID,
	}, nil
}

func NewCreateInformationUsecase(informationRepository repository.InformationRepository) CreateInformationUsecase {
	return &createInformationUsecase{
		informationRepository: informationRepository,
	}
}
