package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/keeper/internal/app/domain"
	"github.com/inter-hubly/pilot/database/pgsql"
	"github.com/inter-hubly/pilot/hlog"
)

type Client interface {
	GetClientById(ctx context.Context, clientId string) (*domain.Client, error)
}

type clientRepository struct {
	connection pgsql.SqlConn
}

func NewClient() *clientRepository {

	var (
		clientRepositoryOnce sync.Once
		client               *clientRepository
	)

	clientRepositoryOnce.Do(func() {
		client = &clientRepository{
			connection: pgsql.GetConnection(),
		}
	})
	return client
}

func (c *clientRepository) GetClientById(ctx context.Context, clientId string) (*domain.Client, error) {
	query := `SELECT c.id, c.name, c.email, c.app_id, c.phone_number_id, c.business_id, c.access_token 
          FROM clients c 
          WHERE c.id = $1`

	queryExec, err := c.connection.Query(query, clientId)
	if err != nil {
		hlog.Error("clientRepository.GetClientById", fmt.Sprintf("error find clientId %s : %s", clientId, err))
		return nil, err
	}
	var clientDb domain.Client
	if err = queryExec.Scan(
		&clientDb.Id,
		&clientDb.Name,
		&clientDb.Email,
		&clientDb.AppId,
		&clientDb.PhoneNumberId,
		&clientDb.BusinessId,
		&clientDb.AccessToken,
	); err != nil {
		hlog.Error("clientRepository.GetClientById", fmt.Sprintf("error scan clientId %s : %s", clientId, err))
		return nil, err
	}
	return &clientDb, nil
}
