package usecase

import (
	"echo-household-budget/internal/infrastructure/persistence/repository"
	"echo-household-budget/internal/model"
)

type WebsocketUsecase interface {
	CreateKaimemo(req model.TelegraphRequest) (interface{}, error)
	RemoveKaimemo(req model.TelegraphRequest) (interface{}, error)
}

type websocketUsecase struct {
	repo repository.KaimemoRepository
}
