package usecase

import (
	domainmodel "echo-household-budget/internal/domain/model"
	"echo-household-budget/internal/domain/repository"
)

type (
	FetchUserInformationInput struct {
		UserID int
	}

	FetchUserInformationOutput struct {
		ID       int
		Title    string
		Content  string
		IsRead   bool
		Category string
	}

	fetchUserInformationUsecase struct {
		userInformationRepository repository.UserInformationRepository
	}

	FetchUserInformationUsecase interface {
		Execute(input FetchUserInformationInput) ([]FetchUserInformationOutput, error)
	}
)

// Execute implements FetchUserInformationUsecase.
func (f *fetchUserInformationUsecase) Execute(input FetchUserInformationInput) ([]FetchUserInformationOutput, error) {
	userInformationModels, err := f.userInformationRepository.FindAllIsPublished(input.UserID)
	if err != nil {
		return nil, err
	}

	output := f.makeOutput(userInformationModels)
	return output, nil
}

func (f *fetchUserInformationUsecase) makeOutput(informations []*domainmodel.UserInformation) []FetchUserInformationOutput {
	output := make([]FetchUserInformationOutput, len(informations))
	for i, information := range informations {
		output[i] = FetchUserInformationOutput{
			ID:       information.InformationID,
			Title:    information.Information.Title,
			Content:  information.Information.Content,
			Category: information.Information.Category,
			IsRead:   information.IsRead,
		}
	}
	return output
}

func NewFetchUserInformationUsecase(userInformationRepository repository.UserInformationRepository) FetchUserInformationUsecase {
	return &fetchUserInformationUsecase{
		userInformationRepository: userInformationRepository,
	}
}
