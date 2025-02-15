package mediator

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/keeper/internal/app/domain"
	"github.com/inter-hubly/keeper/internal/app/domain/kdto"
	"github.com/inter-hubly/keeper/internal/app/gateway"
	"github.com/inter-hubly/keeper/internal/app/repository"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
)

type Template interface {
	Save(ctx context.Context, user *hctx.Logged, dto *kdto.Template) (*domain.Template, error)
}

var (
	templateMediatorOnce sync.Once
	template             *templateMediator
)

type templateMediator struct {
	templateRepository repository.Template
	whatsAppGateway    gateway.WhatsApp
}

func NewTemplate(ctx context.Context) *templateMediator {
	templateMediatorOnce.Do(func() {
		template = &templateMediator{
			templateRepository: repository.NewTemplate(ctx),
			whatsAppGateway:    gateway.NewWhatsApp(ctx),
		}
	})
	return template
}

func (t *templateMediator) Save(ctx context.Context, user *hctx.Logged, templateDomain *kdto.Template) (*domain.Template, error) {
	hlog.Debug(ctx, "templateMediator.Save", fmt.Sprintf("%s", templateDomain))

	t.whatsAppGateway.CreateTemplate(ctx, templateDomain)

	saveTemplate, err := t.templateRepository.SaveTemplate(ctx, user, templateDomain)
	if err != nil {
		hlog.Error(ctx, "templateMediator.Save", fmt.Sprintf("%s", err))
		return nil, err
	}
	return saveTemplate, nil
}
