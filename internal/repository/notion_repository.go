//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/mock_$GOFILE -package=mock
package repository

import (
	"context"
	"log"
	"template-echo-notion-integration/internal/model"

	"github.com/jomei/notionapi"
)

type kaimemoRepository struct {
	client                         *notionapi.Client
	query                          *notionapi.DatabaseQueryRequest
	databaseKaimemoInputID         string
	databaseKaimemoSummaryRecordID string
}

// FetchKaimemoAmount implements KaimemoRepository.
func (k *kaimemoRepository) FetchKaimemoAmountRecords(userID string) (*model.KaimemoAmountRecords, error) {
	k.query.Filter = &notionapi.PropertyFilter{
		Property: "tempUserID",
		RichText: &notionapi.TextFilterCondition{
			Contains: userID,
		},
	}

	resp, err := k.client.Database.Query(context.Background(), notionapi.DatabaseID(k.databaseKaimemoSummaryRecordID), k.query)
	if err != nil {
		log.Printf("failed to notion query database: %v", err)
		return nil, err
	}

	var kaimemoAmounts []model.KaimemoAmount
	for _, result := range resp.Results {
		properties := result.Properties

		data := model.KaimemoAmount{}
		data.ID = string(result.ID)
		for _, property := range properties {
			switch prop := property.(type) {
			case *notionapi.TitleProperty:
				for _, text := range prop.Title {
					data.Date = text.Text.Content
				}
			case *notionapi.NumberProperty:
				data.Amount = int(prop.Number)
			case *notionapi.SelectProperty:
				data.Tag = prop.Select.Name
			default:
				// Unhandled property type
			}
		}
		kaimemoAmounts = append(kaimemoAmounts, data)
	}

	return &model.KaimemoAmountRecords{
		Records: kaimemoAmounts,
	}, nil
}

// InsertKaimemoAmount implements KaimemoRepository.
func (k *kaimemoRepository) InsertKaimemoAmount(req model.CreateKaimemoAmountRequest) error {
	_, err := k.client.Page.Create(context.Background(), &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: notionapi.DatabaseID(k.databaseKaimemoSummaryRecordID),
		},
		Properties: notionapi.Properties{
			"tempUserID": &notionapi.RichTextProperty{
				RichText: []notionapi.RichText{
					{
						Text: &notionapi.Text{
							Content: req.TempUserID,
						},
					},
				},
			},
			"date": &notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{
						Text: &notionapi.Text{
							Content: req.Date,
						},
					},
				},
			},
			"tag": &notionapi.SelectProperty{
				Select: notionapi.Option{
					Name: req.Tag,
				},
			},
			"amount": &notionapi.NumberProperty{
				Number: float64(req.Amount),
			},
		},
	})

	if err != nil {
		log.Printf("failed to notion create page: %v", err)
		return err
	}
	return nil
}

// RemoveKaimemoAmount implements KaimemoRepository.
func (k *kaimemoRepository) RemoveKaimemoAmount(id string, userID string) error {
	_, err := k.client.Page.Update(context.Background(), notionapi.PageID(id), &notionapi.PageUpdateRequest{
		Archived: true,
	})

	if err != nil {
		log.Printf("failed to notion update page: %v", err)
		return err
	}

	return nil
}

// FetchKaimemo implements KaimemoRepository.
func (k *kaimemoRepository) FetchKaimemo(userID string) ([]model.KaimemoResponse, error) {
	k.query.Filter = &notionapi.PropertyFilter{
		Property: "tempUserID",
		RichText: &notionapi.TextFilterCondition{
			Contains: userID,
		},
	}
	resp, err := k.client.Database.Query(context.Background(), notionapi.DatabaseID(k.databaseKaimemoInputID), k.query)
	if err != nil {
		log.Printf("failed to notion query database: %v", err)
		return nil, err
	}

	var kaimemoResponses []model.KaimemoResponse
	for _, result := range resp.Results {
		properties := result.Properties

		data := model.KaimemoResponse{}
		data.ID = string(result.ID)
		for _, property := range properties {
			switch prop := property.(type) {
			case *notionapi.TitleProperty:
				for _, text := range prop.Title {
					data.Name = text.Text.Content
				}
			case *notionapi.SelectProperty:
				data.Tag = prop.Select.Name
			case *notionapi.CheckboxProperty:
				data.Done = prop.Checkbox
			default:
				// fmt.Printf("  %s: Unhandled property type\n", key)
			}
		}
		kaimemoResponses = append(kaimemoResponses, data)
	}

	return kaimemoResponses, nil
}

// InsertKaimemo implements KaimemoRepository.
func (k *kaimemoRepository) InsertKaimemo(req model.CreateKaimemoRequest) error {
	_, err := k.client.Page.Create(context.Background(), &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: notionapi.DatabaseID(k.databaseKaimemoInputID), // 既存のデータベースID
		},
		Properties: notionapi.Properties{
			"tempUserID": &notionapi.RichTextProperty{
				RichText: []notionapi.RichText{
					{
						Text: &notionapi.Text{
							Content: req.TempUserID,
						},
					},
				},
			},
			"name": &notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{
						Text: &notionapi.Text{
							Content: req.Name,
						},
					},
				},
			},
			"tag": &notionapi.SelectProperty{
				Select: notionapi.Option{
					Name: req.Tag,
				},
			},
		},
	})

	if err != nil {
		log.Printf("failed to notion create page: %v", err)
		return err
	}

	return nil
}

// RemoveKaimemo implements KaimemoRepository.
func (k *kaimemoRepository) RemoveKaimemo(id string, userID string) error {
	_, err := k.client.Page.Update(context.Background(), notionapi.PageID(id), &notionapi.PageUpdateRequest{
		Archived: true,
	})

	if err != nil {
		log.Printf("failed to notion update page: %v", err)
		return err
	}

	return nil
}

type KaimemoRepository interface {
	FetchKaimemo(userID string) ([]model.KaimemoResponse, error)
	InsertKaimemo(req model.CreateKaimemoRequest) error
	RemoveKaimemo(id string, userID string) error
	FetchKaimemoAmountRecords(userID string) (*model.KaimemoAmountRecords, error)
	InsertKaimemoAmount(req model.CreateKaimemoAmountRequest) error
	RemoveKaimemoAmount(id string, userID string) error
}

func NewNotionRepository(apiKey string, databaseKaimemoInputID string, databaseKaimemoSummaryRecordID string) KaimemoRepository {
	client := notionapi.NewClient(notionapi.Token(apiKey))
	query := &notionapi.DatabaseQueryRequest{}

	return &kaimemoRepository{client: client, databaseKaimemoInputID: databaseKaimemoInputID, databaseKaimemoSummaryRecordID: databaseKaimemoSummaryRecordID, query: query}
}
