//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go -package=mocks

package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/inter-hubly/pilot/database/pgsql"
	"github.com/inter-hubly/pilot/domain/entity"
	"github.com/inter-hubly/pilot/hlog"
)

type Repository interface {
	GetClientByPhoneNumberId(ctx context.Context, clientId string) (*entity.Client, error)
	SaveClient(ctx context.Context, clientEntity *entity.Client) error
}

type clientRepository struct {
	connection pgsql.SqlConn
}

var (
	_clientRepositoryOnce sync.Once
	_clientRepository     *clientRepository
)

func NewRepository(ctx context.Context) *clientRepository {
	_clientRepositoryOnce.Do(func() {
		_clientRepository = &clientRepository{
			connection: pgsql.GetConnection(ctx),
		}
	})
	return _clientRepository
}

func (c *clientRepository) GetClientByPhoneNumberId(ctx context.Context, clientId string) (*entity.Client, error) {
	hlog.Debug(ctx, "client.repository.GetClientByPhoneNumberId", fmt.Sprintf("geting client id: %s", clientId))
	query := `SELECT c.id, c.name, c.email, c.app_id, c.phone_number_id, c.business_id, c.access_token
          FROM client c 
          WHERE c.phone_number_id = $1`

	queryExec, err := c.connection.Query(query, clientId)
	if err != nil {
		hlog.Error(ctx, "client.repository.GetClientById", fmt.Sprintf("error find clientId %s : %s", clientId, err))
		return nil, err
	}
	var clientDb entity.Client
	if err = queryExec.Scan(
		&clientDb.Id,
		&clientDb.Name,
		&clientDb.Email,
		&clientDb.AppId,
		&clientDb.PhoneNumberId,
		&clientDb.BusinessId,
		&clientDb.AccessToken,
	); err != nil {
		hlog.Error(ctx, "client.repository.GetClientById", fmt.Sprintf("error scan clientId %s : %s", clientId, err))
		return nil, err
	}
	return &clientDb, nil
}

func (c *clientRepository) SaveClient(ctx context.Context, clientEntity *entity.Client) error {
	query := `INSERT INTO client (name, email, app_id, phone_number_id, business_id, access_token, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	row, err := c.connection.Query(query,
		clientEntity.Name,
		clientEntity.Email,
		clientEntity.AppId,
		clientEntity.PhoneNumberId,
		clientEntity.BusinessId,
		clientEntity.AccessToken,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		hlog.Error(ctx, "client.repository.SaveClient", fmt.Sprintf("error save client : %s", err))
		return err
	}
	var returnedId string
	if err = row.Scan(&returnedId); err != nil {
		hlog.Error(ctx, "client.repository.SaveClient", fmt.Sprintf("error save client : %s", err))
		return err
	}
	clientEntity.Id = returnedId
	return nil
}
