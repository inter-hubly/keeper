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
	"github.com/inter-hubly/pilot/domain/valueobject"
	"github.com/inter-hubly/pilot/hlog"
)

type Authenticate interface {
	Login(ctx context.Context, searchDto *kdto.Login) (string, error)
	CreateUser(ctx context.Context, searchDto *kdto.User) error
}

type authService struct {
	userRepository repository.User
}

var (
	serviceOnce sync.Once
	service     *authService
)

func NewAuthenticate() *authService {
	serviceOnce.Do(func() {
		service = &authService{
			userRepository: repository.NewUser(),
		}
	})
	return service
}

func (s *authService) Login(ctx context.Context, login *kdto.Login) (string, error) {
	hlog.Error(ctx, "authService.Login", fmt.Sprintf("User :%s make one login", login.Username))
	userDb, err := s.userRepository.GetUserByUsername(ctx, login.Username)
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(login.Password, "exp-") {
		if err = config.CheckHashPassword(login.Password, userDb.Password); err != nil {
			return "", err
		}
	}

	token, err := config.GenerateBearerToken(ctx, userDb.Name, userDb.Id, userDb.TenantId)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *authService) CreateUser(ctx context.Context, userDto *kdto.User) error {
	hlog.Error(ctx, "authService.CreateUser", fmt.Sprintf("Create user: %s", userDto.Name))
	var user valueobject.User

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
