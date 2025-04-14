package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/keeper/internal/app/domain"
	"github.com/inter-hubly/keeper/internal/app/domain/kdto"
	"github.com/inter-hubly/keeper/internal/app/repository"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
)

type Messages interface {
	SearchMessages(ctx context.Context, loggedUser *hctx.Logged, searchDto *kdto.Search) (map[string]*domain.Conversations, error)
}

type messagesService struct {
	messagesRepository repository.Messages
	contactRepository  repository.Contact
	campaignRepository repository.Campaign
}

var (
	messageServiceOnce sync.Once
	messages           *messagesService
)

func NewMessages(ctx context.Context) *messagesService {
	messageServiceOnce.Do(func() {
		messages = &messagesService{
			messagesRepository: repository.NewMessages(ctx),
			contactRepository:  repository.NewContact(ctx),
			campaignRepository: repository.NewCampaign(ctx),
		}
	})
	return messages
}

// SearchMessages TODO melhorar todo o search
func (s *messagesService) SearchMessages(ctx context.Context, loggedUser *hctx.Logged, searchDto *kdto.Search) (map[string]*domain.Conversations, error) {
	hlog.Debug(ctx, "messagesService.SearchMessages", fmt.Sprintf("search message for logged User: %s", loggedUser))
	msgDb, campaingIds, err := s.messagesRepository.GetMessagesByClientId(ctx, loggedUser.Tenant, searchDto)
	if err != nil {
		hlog.Error(ctx, "messagesService.SearchMessages", fmt.Sprintf("error find messages %s", err))
		return nil, err
	}
	campaings, err := s.campaignRepository.GetCampaignsByIds(ctx, campaingIds...)
	if err != nil {
		hlog.Error(ctx, "messagesService.SearchMessages", fmt.Sprintf("error find campaigns %s", err))
		return nil, err
	}

	templateIdAndMessages := make(map[string]string)
	for _, camp := range campaings {
		templateIdAndMessages[camp.Id] = camp.Template.Message
	}

	if msgDb != nil {
		contacts, err := s.contactRepository.FindContacts(ctx, loggedUser.Tenant)
		if err != nil {
			hlog.Error(ctx, "messagesService.SearchMessages", fmt.Sprintf("error find contacts %s", err))
		}

		for _, ct := range contacts {
			if entry, exists := msgDb[ct.Phone]; exists {
				entry.LocalProfileName = ct.Name
			}
		}

		// TODO melhorar urgente
		for _, value := range msgDb {
			for i := range value.Messages {
				msg := &value.Messages[i]
				if msg.Text != nil {
					if entry, exists := templateIdAndMessages[*msg.Text]; exists {
						msg.Text = &entry
					}
				}
			}
		}

		return msgDb, nil
	}
	return map[string]*domain.Conversations{}, nil
}
