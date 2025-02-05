package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/keeper/internal/app/domain/kdto"
	"github.com/inter-hubly/keeper/internal/app/repository"
	"github.com/inter-hubly/pilot/domain/base"
	"github.com/inter-hubly/pilot/domain/entity"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
)

type Client interface {
	GetClientByPhoneNumberId(ctx context.Context, clientID string) (*entity.Client, error)
	SaveClient(ctx context.Context, user *hctx.Logged, dto kdto.Client) error
}

type clientService struct {
	clientRepository repository.Client
	userRepository   repository.User
}

var (
	clientServiceOnce sync.Once
	client            *clientService
)

func NewClient(ctx context.Context) *clientService {
	clientServiceOnce.Do(func() {
		client = &clientService{
			clientRepository: repository.NewClient(ctx),
			userRepository:   repository.NewUser(ctx),
		}
	})
	return client
}

func (c *clientService) GetClientByPhoneNumberId(ctx context.Context, phoneNumberId string) (*entity.Client, error) {
	clientDb, err := c.clientRepository.GetClientByPhoneNumberId(ctx, phoneNumberId)
	if err != nil {
		hlog.Error(ctx, "clientService.GetClient", fmt.Sprintf("error getting client :%d", err))
		return nil, err
	}
	return clientDb, nil
}

func (c *clientService) SaveClient(ctx context.Context, loggedUser *hctx.Logged, clientDto kdto.Client) error {
	clientEntity := entity.Client{
		Name:          clientDto.Name,
		Email:         clientDto.Email,
		AppId:         clientDto.AppId,
		PhoneNumberId: clientDto.PhoneNumberId,
		BusinessId:    clientDto.BusinessId,
		AccessToken:   clientDto.AccessToken,
		Entity:        base.NewBaseEntity(ctx, loggedUser),
	}
	if err := c.clientRepository.SaveClient(ctx, &clientEntity); err != nil {
		hlog.Error(ctx, "clientService.SaveClient", fmt.Sprintf("error saving client :%s", err))
		return err
	}

	if err := c.userRepository.UpdateClientInUser(ctx, clientEntity.PhoneNumberId, loggedUser.UserId); err != nil {
		hlog.Error(ctx, "clientService.SaveClient", fmt.Sprintf("error saving client :%s", err))
		return err
	}
	return nil
}
