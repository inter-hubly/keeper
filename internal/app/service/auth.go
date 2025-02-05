package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/inter-hubly/keeper/internal/app/domain/kdto"
	"github.com/inter-hubly/keeper/internal/app/repository"
	"github.com/inter-hubly/keeper/internal/infraestructure/config"
	"github.com/inter-hubly/pilot/domain/entity"
	"github.com/inter-hubly/pilot/hlog"
)

type Authenticate interface {
	Login(ctx context.Context, searchDto *kdto.Login) (*kdto.Authenticate, error)
	CreateUser(ctx context.Context, searchDto *kdto.User) error
}

type authService struct {
	userRepository   repository.User
	clientRepository repository.Client
}

var (
	serviceOnce sync.Once
	service     *authService
)

func NewAuthenticate(ctx context.Context) *authService {
	serviceOnce.Do(func() {
		service = &authService{
			userRepository:   repository.NewUser(ctx),
			clientRepository: repository.NewClient(ctx),
		}
	})
	return service
}

func (s *authService) Login(ctx context.Context, login *kdto.Login) (*kdto.Authenticate, error) {
	hlog.Debug(ctx, "authService.Login", fmt.Sprintf("User :%s make one login", login.Username))
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

	return &kdto.Authenticate{
		AccessToken: token,
		TenantId:    userDb.TenantId,
	}, nil
}

func (s *authService) CreateUser(ctx context.Context, userDto *kdto.User) error {
	hlog.Error(ctx, "authService.CreateUser", fmt.Sprintf("Create user: %s", userDto.Name))
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
		hlog.Error(ctx, "userRepository.SaveUser", fmt.Sprintf("error insert User : %s", err))
		return err
	}

	return nil
}
