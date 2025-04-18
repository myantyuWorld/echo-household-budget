//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/mock_$GOFILE -package=mock
package usecase

import (
	"echo-household-budget/internal/infrastructure/persistence/repository"
	"echo-household-budget/internal/model"
)

type KaimemoService interface {
	FetchKaimemo(userID string) ([]model.KaimemoResponse, error)
	CreateKaimemo(req model.CreateKaimemoRequest) error
	RemoveKaimemo(id string, userID string) error
	FetchKaimemoSummaryRecord(userID string) (model.KaimemoSummaryResponse, error)
	CreateKaimemoAmount(req model.CreateKaimemoAmountRequest) error
	RemoveKaimemoAmount(id string, userID string) error
}

func NewKaimemoService(repo repository.KaimemoRepository) KaimemoService {
	return &kaimemoService{repo: repo}
}

type kaimemoService struct {
	repo repository.KaimemoRepository
	// TODO : shoppingRepositoryに切り替え（Notion >>> RDBMS）
}

// CreateKaimemoAmount implements KaimemoService.
func (k *kaimemoService) CreateKaimemoAmount(req model.CreateKaimemoAmountRequest) error {
	return k.repo.InsertKaimemoAmount(req)
}

// FetchKaimemoSummaryRecord implements KaimemoService.
func (k *kaimemoService) FetchKaimemoSummaryRecord(userID string) (model.KaimemoSummaryResponse, error) {
	res, err := k.repo.FetchKaimemoAmountRecords(userID)
	if err != nil {
		return model.KaimemoSummaryResponse{
			MonthlySummaries: []model.MonthlySummary{},
			WeeklySummaries:  []model.WeeklySummary{},
		}, err
	}

	return model.KaimemoSummaryResponse{
		MonthlySummaries: res.GroupByMonth(),
		WeeklySummaries:  res.GroupByWeek(),
	}, nil
}

// RemoveKaimemoAmount implements KaimemoService.
func (k *kaimemoService) RemoveKaimemoAmount(id string, userID string) error {
	return k.repo.RemoveKaimemoAmount(id, userID)
}

// CreateKaimemo implements KaimemoService.
func (k *kaimemoService) CreateKaimemo(req model.CreateKaimemoRequest) error {
	return k.repo.InsertKaimemo(req)
}

// FetchKaimemo implements KaimemoService.
func (k *kaimemoService) FetchKaimemo(userID string) ([]model.KaimemoResponse, error) {
	return k.repo.FetchKaimemo(userID)
}

// RemoveKaimemo implements KaimemoService.
func (k *kaimemoService) RemoveKaimemo(id string, userID string) error {
	return k.repo.RemoveKaimemo(id, userID)
}
