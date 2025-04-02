package mediator

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/inter-hubly/keeper/internal/app/domain"
	"github.com/inter-hubly/keeper/internal/app/gateway"
	"github.com/inter-hubly/keeper/internal/app/repository"
	"github.com/inter-hubly/pilot/domain/base"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
	"github.com/inter-hubly/pilot/server"
	"github.com/inter-hubly/pilot/util"
)

type Template interface {
	Save(ctx context.Context, user *hctx.Logged, dto *domain.Template) (*domain.Template, error)
}

var (
	_templateMediatorOnce sync.Once
	_template             *templateMediator
)

type templateMediator struct {
	templateRepository repository.Template
	whatsAppGateway    gateway.WhatsApp
}

func NewTemplate(ctx context.Context) *templateMediator {
	_templateMediatorOnce.Do(func() {
		env := server.GetEnvironment().Env
		var gtw gateway.WhatsApp
		if env == server.Development {
			gtw = gateway.NewWhatsAppMock()
		} else {
			gtw = gateway.NewWhatsApp(ctx)
		}
		_template = &templateMediator{
			templateRepository: repository.NewTemplate(ctx),
			whatsAppGateway:    gtw,
		}
	})
	return _template
}

func (t *templateMediator) Save(ctx context.Context, user *hctx.Logged, templateDomain *domain.Template) (*domain.Template, error) {
	hlog.Debug(ctx, "templateMediator.Save", fmt.Sprint(templateDomain))

	// for each component, header, body and footer
	re := regexp.MustCompile(`\{\{([^}]+)\}\}`)
	for _, v := range templateDomain.Components {
		matches := re.FindAllStringSubmatch(v.Text, -1)

		if matches != nil {
			manyExamples := v.Example[strings.ToLower(fmt.Sprintf("%s_text", v.Type))]

			// many examples
			for _, eachExample := range manyExamples {
				// need verify each example has the same value then match
				if len(eachExample) != len(matches) {
					return nil, errors.New("error when save template")
				}
			}
		}

	}

	templateDomain.Slug = util.ToSlug(templateDomain.Name, true)
	gatewayResponse, err := t.whatsAppGateway.CreateTemplate(ctx, templateDomain)

	if err != nil {
		hlog.Error(ctx, "templateMediator.Save", err.Error())
		return nil, err
	}
	templateDomain.Entity = base.NewBaseEntity(ctx, user)
	templateDomain.ResponseId = gatewayResponse.Id
	templateDomain.Status = domain.TemplateStatus(gatewayResponse.Status)

	saveTemplate, err := t.templateRepository.SaveTemplate(ctx, user, templateDomain)
	if err != nil {
		hlog.Error(ctx, "templateMediator.Save", fmt.Sprintf("%s", err))
		return nil, err
	}
	return saveTemplate, nil
}
