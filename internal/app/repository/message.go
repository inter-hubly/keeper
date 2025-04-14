package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/keeper/internal/app/domain"
	"github.com/inter-hubly/keeper/internal/app/domain/kdto"
	"github.com/inter-hubly/pilot/database/elasticsearch"
	"github.com/inter-hubly/pilot/hlog"
)

var (
	_messageRepositoryOnce sync.Once
	_messageRepository     *messagesRepository
)

type Messages interface {
	GetMessagesByClientId(ctx context.Context, clientId string, searchDto *kdto.Search) (map[string]*domain.Conversations, []string, error)
}

type messagesRepository struct {
	elasticIndex string
	connection   elasticsearch.ElasticConn
}

func NewMessages(ctx context.Context) *messagesRepository {
	_messageRepositoryOnce.Do(func() {
		_messageRepository = &messagesRepository{
			elasticIndex: "whatsapp",
			connection:   elasticsearch.GetConnection(),
		}
	})

	return _messageRepository
}

func (m *messagesRepository) GetMessagesByClientId(ctx context.Context, ownerId string, searchDto *kdto.Search) (map[string]*domain.Conversations, []string, error) {
	fields := map[string]interface{}{
		"size": 100,
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"sort": []interface{}{
			map[string]interface{}{
				"createdAt": map[string]interface{}{
					"order": "asc",
				},
			},
		},
	}

	allMessages, err := m.connection.FindAll(ctx, fmt.Sprintf("%s.%s", ownerId, m.elasticIndex), fields)
	if err != nil {
		hlog.Error(ctx, "messagesRepository.GetMessagesByClientId", err.Error())
		return nil, nil, err
	}
	if allMessages == nil || allMessages.Hits == nil {
		return nil, nil, nil
	}

	conversations := make(map[string]*domain.Conversations)
	manyCampaigns := make([]string, 0, len(allMessages.Hits.Hits))

	var (
		foundProfileName bool
		conv             *domain.Conversations
	)

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
				manyCampaigns = append(manyCampaigns, val["campaignId"].(string))
				message = val["campaignId"].(string)
			}
			var messageId string
			if msgIdValue, msgIdok := val["messageId"].(string); msgIdok {
				messageId = msgIdValue
			}

			singleMessage := domain.Message{
				Id:      messageId,
				IsOwner: val["isOwner"].(bool),
				Text:    &message,
			}
			conv.Messages = append(conv.Messages, singleMessage)
			if !foundProfileName {
				if profileNameString, ok := val["profileName"].(string); ok {
					conv.WhatsAppProfileName = profileNameString
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

	return conversations, manyCampaigns, nil
}
