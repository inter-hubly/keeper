package service

import (
	"context"
	"sync"
)

type Client interface {
	GetClient(ctx context.Context, clientID string) (Client, error)
}

var (
	clientServiceOnce sync.Once
	client            *clientService
)

type clientService struct {
}

func NewClient() *clientService {
	clientServiceOnce.Do(func() {
		client = &clientService{}
	})
	return client
}

func (c *clientService) GetClient(ctx context.Context, clientID string) (Client, error) {
	return nil, nil
}
