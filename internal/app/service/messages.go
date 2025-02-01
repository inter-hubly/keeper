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
}

var (
	messageServiceOnce sync.Once
	messages           *messagesService
)

func NewMessages() *messagesService {
	messageServiceOnce.Do(func() {
		messages = &messagesService{
			messagesRepository: repository.NewMessages(),
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
		return msgDb, nil
	}
	return map[string]*domain.Conversations{}, nil
}
