package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/keeper/internal/app/repository"
	"github.com/inter-hubly/pilot/domain/valueobject"
	"github.com/inter-hubly/pilot/hlog"
)

type Client interface {
	GetClientByPhoneNumberId(ctx context.Context, clientID string) (*valueobject.Client, error)
}

type clientService struct {
	clientRepository repository.Client
}

var (
	clientServiceOnce sync.Once
	client            *clientService
)

func NewClient() *clientService {
	clientServiceOnce.Do(func() {
		client = &clientService{
			clientRepository: repository.NewClient(),
		}
	})
	return client
}

func (c *clientService) GetClientByPhoneNumberId(ctx context.Context, phoneNumberId string) (*valueobject.Client, error) {
	clientDb, err := c.clientRepository.GetClientByPhoneNumberId(ctx, phoneNumberId)
	if err != nil {
		hlog.Error(ctx, "clientService.GetClient", fmt.Sprintf("error getting client :%d", err))
		return nil, err
	}
	return clientDb, nil
}
