//go:generate mockgen -source=user.go -destination=mocks/user_mock.go -package=mocks

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/inter-hubly/pilot/database/pgsql"
	"github.com/inter-hubly/pilot/domain/entity"
	"github.com/inter-hubly/pilot/hlog"
)

type User interface {
	GetUserByUsername(ctx context.Context, emailUser string) (*entity.User, error)
	SaveUser(ctx context.Context, user *entity.User) error
	UpdateClientInUser(ctx context.Context, clientId, userId string) error
}

type userRepository struct {
	connection pgsql.SqlConn
}

var (
	_userRepositoryOnce sync.Once
	_userRepository     *userRepository
)

func NewUser(ctx context.Context) *userRepository {

	_userRepositoryOnce.Do(func() {
		_userRepository = &userRepository{
			connection: pgsql.GetConnection(ctx),
		}
	})
	return _userRepository
}

func (r *userRepository) GetUserByUsername(ctx context.Context, userEmail string) (*entity.User, error) {
	hlog.Debug(ctx, "userRepository.GetUserByUsername", userEmail)
	query := `SELECT u.id, u.name, u.email, u.password, u.login_attempt, u.tenant_id, u.created_at, u.updated_at
          FROM "user" u WHERE u.email = $1`

	queryExec, err := r.connection.Query(query, userEmail)
	if err != nil {
		hlog.Error(ctx, "userRepository.GetUserByUsername", fmt.Sprintf("error find UserEmail %s : %s", userEmail, err))
		return nil, err
	}
	var userDb entity.User
	var tenantID sql.NullString
	if err = queryExec.Scan(
		&userDb.Id,
		&userDb.Name,
		&userDb.Email,
		&userDb.Password,
		&userDb.LoginAttempt,
		&tenantID,
		&userDb.CreatedAt,
		&userDb.UpdatedAt,
	); err != nil {
		hlog.Error(ctx, "userRepository.GetUserByUsername", fmt.Sprintf("error scan UserEmail %s : %s", userEmail, err))
		return nil, err
	}
	if tenantID.Valid {
		userDb.TenantId = tenantID.String
	}
	return &userDb, nil
}

func (r *userRepository) SaveUser(ctx context.Context, user *entity.User) error {
	hlog.Debug(ctx, "userRepository.UpdateClientInUser", fmt.Sprintf("save User %+v", user))
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

func (r *userRepository) UpdateClientInUser(ctx context.Context, tenantId, userId string) error {
	hlog.Debug(ctx, "userRepository.UpdateClientInUser", fmt.Sprintf(" Set tenantId: %s and UserId: %s", tenantId, userId))
	query := `UPDATE "user" SET tenant_id = $1 WHERE id = $2`
	if _, err := r.connection.Exec(query, tenantId, userId); err != nil {
		hlog.Error(ctx, "userRepository.UpdateClientInUser", fmt.Sprintf("update User Error : %s", err))
		return err
	}
	return nil
}
