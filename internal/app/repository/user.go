package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/pilot/database/pgsql"
	"github.com/inter-hubly/pilot/domain/valueobject"
	"github.com/inter-hubly/pilot/hlog"
)

type User interface {
	GetUserByUsername(ctx context.Context, clientId string) (*valueobject.User, error)
}

type userRepository struct {
	connection pgsql.SqlConn
}

func NewUser() *userRepository {

	var (
		userOnce   sync.Once
		repository *userRepository
	)

	userOnce.Do(func() {
		repository = &userRepository{
			connection: pgsql.GetConnection(),
		}
	})
	return repository
}

func (r *userRepository) GetUserByUsername(ctx context.Context, userEmail string) (*valueobject.User, error) {
	query := `SELECT c.id, c.name, c.email, c.password, c.client_id, c.login_attempt, c.created_at, c.updated_at 
          FROM "user" c 
          WHERE c.email = $1`

	queryExec, err := r.connection.Query(query, userEmail)
	if err != nil {
		hlog.Error("userRepository.GetUserByUsername", fmt.Sprintf("error find UserEmail %s : %s", userEmail, err))
		return nil, err
	}
	var userDb valueobject.User
	if err = queryExec.Scan(
		&userDb.Id,
		&userDb.Name,
		&userDb.Email,
		&userDb.Password,
		&userDb.ClientId,
		&userDb.LoginAttempt,
		&userDb.CreatedAt,
		&userDb.UpdatedAt,
	); err != nil {
		hlog.Error("userRepository.GetUserByUsername", fmt.Sprintf("error scan UserEmail %s : %s", userEmail, err))
		return nil, err
	}
	return &userDb, nil
}
