package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/keeper/internal/app/domain"
	"github.com/inter-hubly/keeper/internal/app/repository"
	"github.com/inter-hubly/pilot/hlog"
)

type Client interface {
	GetClient(ctx context.Context, clientID string) (*domain.Client, error)
}

type clientService struct {
	clientRepository repository.Client
}

func NewClient() *clientService {

	var (
		clientServiceOnce sync.Once
		client            *clientService
	)

	clientServiceOnce.Do(func() {
		client = &clientService{
			clientRepository: repository.NewClient(),
		}
	})
	return client
}

func (c *clientService) GetClient(ctx context.Context, clientID string) (*domain.Client, error) {
	clientDb, err := c.clientRepository.GetClientById(ctx, clientID)
	if err != nil {
		hlog.Error("clientService.GetClient", fmt.Sprintf("error getting client :%d", err))
		return nil, err
	}
	return clientDb, nil
}
