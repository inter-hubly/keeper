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
	GetUserByUsername(ctx context.Context, emailUser string) (*valueobject.User, error)
	SaveUser(ctx context.Context, user *valueobject.User) error
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
	query := `SELECT u.id, u.name, u.email, u.password, u.client_id, u.login_attempt, u.created_at, u.updated_at, c.phone_number_Id 
          FROM "user" u left join client c on c.id = u.client_id
          WHERE u.email = $1`

	queryExec, err := r.connection.Query(query, userEmail)
	if err != nil {
		hlog.Error(ctx, "userRepository.GetUserByUsername", fmt.Sprintf("error find UserEmail %s : %s", userEmail, err))
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
		&userDb.TenantId,
	); err != nil {
		hlog.Error(ctx, "userRepository.GetUserByUsername", fmt.Sprintf("error scan UserEmail %s : %s", userEmail, err))
		return nil, err
	}
	return &userDb, nil
}

func (r *userRepository) SaveUser(ctx context.Context, user *valueobject.User) error {
	query := `
		INSERT INTO "user" (name, email, password, login_attempt, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6)`
	if _, err := r.connection.Exec(
		query,
		user.Name,
		user.Email,
		user.Password,
		user.LoginAttempt,
		user.CreatedAt,
		user.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}
