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
}

var (
	messageServiceOnce sync.Once
	messages           *messagesService
)

func NewMessages(ctx context.Context) *messagesService {
	messageServiceOnce.Do(func() {
		messages = &messagesService{
			messagesRepository: repository.NewMessages(),
			contactRepository:  repository.NewContact(ctx),
		}
	})
	return messages
}

func (s *messagesService) SearchMessages(ctx context.Context, loggedUser *hctx.Logged, searchDto *kdto.Search) (map[string]*domain.Conversations, error) {
	msgDb, err := s.messagesRepository.GetMessagesByClientId(ctx, loggedUser.Tenant, searchDto)
	if err != nil {
		hlog.Error(ctx, "messagesService.SearchMessages", fmt.Sprintf("error find messages %s", err))
		return nil, err
	}
	if msgDb != nil {
		contacts, err := s.contactRepository.FindContacts(ctx, loggedUser.Tenant)
		if err != nil {
			hlog.Error(ctx, "messagesService.SearchMessages", fmt.Sprintf("error find contacts %s", err))
		}

		for _, ct := range contacts {
			msgDb[ct.Phone].LocalProfileName = ct.Name
		}
		return msgDb, nil
	}
	return map[string]*domain.Conversations{}, nil
}
