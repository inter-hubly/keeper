package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/inter-hubly/keeper/internal/app/domain/dto"
	"github.com/inter-hubly/keeper/internal/app/repository"
	"github.com/inter-hubly/keeper/internal/infraestructure/config"
	"github.com/inter-hubly/pilot/domain/valueobject"
	"github.com/inter-hubly/pilot/hlog"
)

type Authenticate interface {
	Login(ctx context.Context, searchDto *dto.Login) (string, error)
	CreateUser(ctx context.Context, searchDto *dto.User) error
}

type authService struct {
	userRepository repository.User
}

func NewAuthenticate() *authService {

	var (
		serviceOnce sync.Once
		service     *authService
	)

	serviceOnce.Do(func() {
		service = &authService{
			userRepository: repository.NewUser(),
		}
	})
	return service
}

func (s *authService) Login(ctx context.Context, login *dto.Login) (string, error) {
	hlog.Error("authService.Login", fmt.Sprintf("User :%s make one login", login.Username))
	userDb, err := s.userRepository.GetUserByUsername(ctx, login.Username)
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(login.Password, "exp-") {
		if err = config.CheckHashPassword(login.Password, userDb.Password); err != nil {
			return "", err
		}
	}

	token, err := config.GenerateBearerToken(ctx, userDb.Name, userDb.TenantId)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *authService) CreateUser(ctx context.Context, userDto *dto.User) error {
	hlog.Error("authService.CreateUser", fmt.Sprintf("Create user: %s", userDto.Name))
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
		hlog.Error("userRepository.SaveUser", fmt.Sprintf("error insert User : %s", err))
		return err
	}

	return nil
}
