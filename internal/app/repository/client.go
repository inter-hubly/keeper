package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/pilot/database/pgsql"
	"github.com/inter-hubly/pilot/domain/valueobject"
	"github.com/inter-hubly/pilot/hlog"
)

type Client interface {
	GetClientByPhoneNumberId(ctx context.Context, clientId string) (*valueobject.Client, error)
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

func (c *clientRepository) GetClientByPhoneNumberId(ctx context.Context, clientId string) (*valueobject.Client, error) {
	query := `SELECT c.id, c.name, c.email, c.app_id, c.phone_number_id, c.business_id
          FROM client c 
          WHERE c.phone_number_id = $1`

	queryExec, err := c.connection.Query(query, clientId)
	if err != nil {
		hlog.Error(ctx, "clientRepository.GetClientById", fmt.Sprintf("error find clientId %s : %s", clientId, err))
		return nil, err
	}
	var clientDb valueobject.Client
	if err = queryExec.Scan(
		&clientDb.Id,
		&clientDb.Name,
		&clientDb.Email,
		&clientDb.AppId,
		&clientDb.PhoneNumberId,
		&clientDb.BusinessId,
	); err != nil {
		hlog.Error(ctx, "clientRepository.GetClientById", fmt.Sprintf("error scan clientId %s : %s", clientId, err))
		return nil, err
	}
	return &clientDb, nil
}
