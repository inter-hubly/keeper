package repository

import (
	"context"
	"sync"

	"github.com/inter-hubly/keeper/internal/app/domain"
	"github.com/inter-hubly/keeper/internal/app/domain/dto"
	"github.com/inter-hubly/pilot/database/elasticsearch"
	"github.com/inter-hubly/pilot/hlog"
)

type Messages interface {
	GetMessagesByClientId(ctx context.Context, clientId string, searchDto *dto.Search) (map[string]*domain.Conversations, error)
}

type messagesRepository struct {
	elasticIndex string
	connection   elasticsearch.ElasticConn
}

func NewMessages() *messagesRepository {

	var (
		repositoryOnce sync.Once
		repository     *messagesRepository
	)

	repositoryOnce.Do(func() {
		repository = &messagesRepository{
			elasticIndex: "whatsapp.ready",
			connection:   elasticsearch.GetConnection(),
		}
	})
	return repository
}

func (m *messagesRepository) GetMessagesByClientId(ctx context.Context, ownerId string, searchDto *dto.Search) (map[string]*domain.Conversations, error) {
	fields := map[string]interface{}{
		"size": 100,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"match": map[string]interface{}{
							"ownerId": ownerId,
						},
					},
				},
			},
		},
		"sort": []interface{}{
			map[string]interface{}{
				"createdAt": map[string]interface{}{
					"order": "asc",
				},
			},
		},
	}

	allMessages, err := m.connection.FindAll(ctx, m.elasticIndex, fields)
	if err != nil {
		hlog.Error("messagesRepository.GetMessagesByClientId", err.Error())
		return nil, err
	}
	if allMessages == nil || allMessages.Hits == nil {
		return nil, nil
	}
	conversations := make(map[string]*domain.Conversations)
	// response := make([]*domain.Message, 0, len(allMessages.Hits.Hits))

	var foundProfileName bool
	var conv *domain.Conversations

	for _, elasticSingleResponse := range allMessages.Hits.Hits {
		if val, ok := elasticSingleResponse.(map[string]interface{})["_source"].(map[string]interface{}); ok {

			elasticPhoneValue := val["toPhoneId"].(string)
			if existingConv, phoneOk := conversations[elasticPhoneValue]; phoneOk {
				conv = existingConv
			} else {
				foundProfileName = false
				conv = &domain.Conversations{
					Messages: []domain.Message{},
				}
			}
			var message string
			if msgValue, ok := val["message"].(string); ok {
				message = msgValue
			} else {
				message = val["templateName"].(string)
			}
			var messageId string
			if msgIdValue, msgIdok := val["messageId"].(string); msgIdok {
				messageId = msgIdValue
			}

			singleMessage := domain.Message{
				Id:      messageId,
				IsOwner: val["isOwner"].(bool),
				Text:    message,
			}
			conv.Messages = append(conv.Messages, singleMessage)
			if !foundProfileName {
				if profileNameString, ok := val["profileName"].(string); ok {
					conv.ProfileName = profileNameString
					foundProfileName = true
				}
			}

			status := []domain.MessageStatus{}

			if statuses, ok := val["status"]; ok {
				if statusList, valid := statuses.([]interface{}); valid {
					for _, st := range statusList {
						if stMap, isMap := st.(map[string]interface{}); isMap {
							msgStatus := domain.MessageStatus{}

							if receivedAt, exists := stMap["receivedAt"]; exists {
								if receivedAtFloat, isFloat := receivedAt.(float64); isFloat {
									msgStatus.ReceivedAt = receivedAtFloat
								}
							}

							if statusStr, exists := stMap["status"]; exists {
								if statusString, isString := statusStr.(string); isString {
									msgStatus.Status = statusString
								}
							}

							status = append(status, msgStatus)
						}
					}
					singleMessage.Status = status
				}
			}

			conversations[elasticPhoneValue] = conv
		}
	}

	return conversations, nil
}
