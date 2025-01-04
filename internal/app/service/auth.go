package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/keeper/internal/app/domain/dto"
	"github.com/inter-hubly/keeper/internal/app/repository"
	"github.com/inter-hubly/keeper/internal/infraestructure/config"
	"github.com/inter-hubly/pilot/hlog"
)

type Authenticate interface {
	Login(ctx context.Context, searchDto *dto.Login) (string, error)
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
	token, err := config.GenerateBearerToken(ctx, userDb.Name, userDb.ClientId)
	if err != nil {
		return "", err
	}
	return token, nil
}
