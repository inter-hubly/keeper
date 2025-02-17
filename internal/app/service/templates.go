package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/keeper/internal/app/domain"
	"github.com/inter-hubly/keeper/internal/app/repository"
	"github.com/inter-hubly/pilot/hlog"

	"github.com/inter-hubly/pilot/hctx"
)

type Template interface {
	FindAll(ctx context.Context, user *hctx.Logged) ([]domain.Template, error)
}

var (
	templateServiceOnce sync.Once
	template            *templateService
)

type templateService struct {
	templateRepository repository.Template
}

func NewTemplate(ctx context.Context) *templateService {
	templateServiceOnce.Do(func() {
		template = &templateService{
			templateRepository: repository.NewTemplate(ctx),
		}
	})
	return template
}

func (s *templateService) FindAll(ctx context.Context, user *hctx.Logged) ([]domain.Template, error) {
	hlog.Debug(ctx, "templateService.FindAll", fmt.Sprintf("Find All Templates %s", user))
	return s.templateRepository.FindAll(ctx, user)
}
