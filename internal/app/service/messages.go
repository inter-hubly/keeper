package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/keeper/internal/app/domain"
	"github.com/inter-hubly/keeper/internal/app/domain/dto"
	"github.com/inter-hubly/keeper/internal/app/repository"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
)

type Messages interface {
	SearchMessages(ctx context.Context, searchDto *dto.Search) ([]*domain.SingleMessage, error)
}

type messagesService struct {
	messagesRepository repository.Messages
}

func NewMessages() *messagesService {

	var (
		serviceOnce sync.Once
		service     *messagesService
	)

	serviceOnce.Do(func() {
		service = &messagesService{
			messagesRepository: repository.NewMessages(),
		}
	})
	return service
}

func (s *messagesService) SearchMessages(ctx context.Context, searchDto *dto.Search) ([]*domain.SingleMessage, error) {
	tenant := hctx.Tenant.Get(ctx)

	messages, err := s.messagesRepository.GetMessagesByClientId(ctx, tenant, searchDto)
	if err != nil {
		hlog.Error("messagesService.SearchMessages", fmt.Sprintf("error find messages %s", err))
		return nil, err
	}
	return messages, nil
}
