package repository

import domainmodel "echo-household-budget/internal/domain/model"

type (
	ChatMessageRepository interface {
		Create(input *domainmodel.ChatMessage) error
		FindByHouseholdID(householdID int, limit int, offset int) ([]*domainmodel.ChatMessage, error)
	}
)
