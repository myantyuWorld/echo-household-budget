package usecase

import (
	domainmodel "echo-household-budget/internal/domain/model"
	repository "echo-household-budget/internal/domain/repository"

	"github.com/davecgh/go-spew/spew"
)

type (
	FetchInformationInput struct {
		// TODO: 拡張性のため、定義だけしておく
	}

	FetchInformationOutput struct {
		ID          int
		Title       string
		Content     string
		IsPublished bool
		Category    string
	}

	FetchInformationUsecase interface {
		Execute() ([]FetchInformationOutput, error)
	}

	fetchInformationUsecase struct {
		informationRepository repository.InformationRepository
	}
)

func NewFetchInformationUsecase(informationRepository repository.InformationRepository) FetchInformationUsecase {
	return &fetchInformationUsecase{
		informationRepository: informationRepository,
	}
}

func (f *fetchInformationUsecase) Execute() ([]FetchInformationOutput, error) {
	informations, err := f.informationRepository.FindAll()
	if err != nil {
		return nil, err
	}

	spew.Dump(informations)

	output := f.makeOutput(informations)

	return output, nil
}

func (f *fetchInformationUsecase) makeOutput(informations []*domainmodel.Information) []FetchInformationOutput {
	output := make([]FetchInformationOutput, len(informations))
	for i, information := range informations {
		output[i] = FetchInformationOutput{
			ID:          information.ID,
			Title:       information.Title,
			Content:     information.Content,
			IsPublished: information.IsPublished,
			Category:    information.Category,
		}
	}
	return output
}
