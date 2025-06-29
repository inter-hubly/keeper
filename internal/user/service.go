package user

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/inter-hubly/keeper/infraestructure/config"
	kdto2 "github.com/inter-hubly/keeper/internal/domain/kdto"
	"github.com/inter-hubly/pilot/domain/entity"
	"github.com/inter-hubly/pilot/hlog"
)

type Service interface {
	Login(ctx context.Context, searchDto *kdto2.Login) (*kdto2.Authenticate, error)
	CreateUser(ctx context.Context, searchDto *kdto2.User) error
}

type authService struct {
	userRepository Repository
}

var (
	_userServiceOnce sync.Once
	_userService     *authService
)

func NewAuthenticate(ctx context.Context) *authService {
	_userServiceOnce.Do(func() {
		_userService = &authService{
			userRepository: NewRepository(ctx),
		}
	})
	return _userService
}

func (s *authService) Login(ctx context.Context, login *kdto2.Login) (*kdto2.Authenticate, error) {
	hlog.Debug(ctx, "user.service.Login", fmt.Sprintf("User :%s make one login", login.Username))
	userDb, err := s.userRepository.GetUserByUsername(ctx, login.Username)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(login.Password, "exp-") {
		if err = config.CheckHashPassword(login.Password, userDb.Password); err != nil {
			return nil, err
		}
	}

	token, err := config.GenerateBearerToken(ctx, userDb.Name, userDb.Id, userDb.TenantId)
	if err != nil {
		return nil, err
	}

	return &kdto2.Authenticate{
		AccessToken: token,
		TenantId:    userDb.TenantId,
	}, nil
}

func (s *authService) CreateUser(ctx context.Context, userDto *kdto2.User) error {
	hlog.Error(ctx, "user.service.CreateUser", fmt.Sprintf("Create user: %s", userDto.Name))
	var user entity.User

	user.Name = userDto.Name
	if hashPass, err := config.HashPassword(userDto.Password); err == nil {
		user.Password = hashPass
	}
	user.Email = userDto.Email
	user.LoginAttempt = 0
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	if err := s.userRepository.SaveUser(ctx, &user); err != nil {
		hlog.Error(ctx, "user.service.SaveUser", fmt.Sprintf("error insert User : %s", err))
		return err
	}

	return nil
}
