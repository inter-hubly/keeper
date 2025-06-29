package messages

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/keeper/internal/campaign"
	"github.com/inter-hubly/keeper/internal/contact"
	"github.com/inter-hubly/keeper/internal/domain"
	"github.com/inter-hubly/keeper/internal/domain/kdto"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
)

type Service interface {
	SearchMessages(ctx context.Context, loggedUser *hctx.Logged, searchDto *kdto.Search) (map[string]*domain.Conversations, error)
}

type messagesService struct {
	messagesRepository Repository
	contactRepository  contact.Repository
	campaignRepository campaign.Repository
}

var (
	_messageServiceOnce sync.Once
	_messageService     *messagesService
)

func NewService(ctx context.Context) *messagesService {
	_messageServiceOnce.Do(func() {
		_messageService = &messagesService{
			messagesRepository: NewRepository(ctx),
			contactRepository:  contact.NewContact(ctx),
			campaignRepository: campaign.NewRepository(ctx),
		}
	})
	return _messageService
}

// SearchMessages TODO melhorar todo o search
func (s *messagesService) SearchMessages(ctx context.Context, loggedUser *hctx.Logged, searchDto *kdto.Search) (map[string]*domain.Conversations, error) {
	hlog.Debug(ctx, "messages.service.SearchMessages", fmt.Sprintf("search message for logged User: %s", loggedUser))
	msgDb, campaingIds, err := s.messagesRepository.GetMessagesByClientId(ctx, loggedUser.Tenant, searchDto)
	if err != nil {
		hlog.Error(ctx, "messages.service.SearchMessages", fmt.Sprintf("error find messages %s", err))
		return nil, err
	}

	templateIdAndMessages := make(map[string]string)
	if len(campaingIds) > 0 {
		campaings, err := s.campaignRepository.GetCampaignsByIds(ctx, campaingIds...)
		if err != nil {
			hlog.Error(ctx, "messages.service.SearchMessages", fmt.Sprintf("error find campaigns %s", err))
			return nil, err
		}
		for _, camp := range campaings {
			templateIdAndMessages[camp.Id] = camp.Template.Message
		}

	}

	if msgDb != nil {
		contacts, err := s.contactRepository.FindContacts(ctx, loggedUser.Tenant)
		if err != nil {
			hlog.Error(ctx, "messages.service.SearchMessages", fmt.Sprintf("error find contacts %s", err))
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
