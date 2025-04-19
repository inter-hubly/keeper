package mediator

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
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
	variablesRepository repository.Variable
	templateRepository  repository.Template
	whatsAppGateway     gateway.WhatsApp
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
			variablesRepository: repository.NewVariables(ctx),
			templateRepository:  repository.NewTemplate(ctx),
			whatsAppGateway:     gtw,
		}
	})
	return _template
}

func (t *templateMediator) Save(ctx context.Context, user *hctx.Logged, templateDomain *domain.Template) (*domain.Template, error) {
	hlog.Debug(ctx, "templateMediator.Save", fmt.Sprint(templateDomain))

	cloned := make([]domain.Components, len(templateDomain.Components))
	copy(cloned, templateDomain.Components)

	variables, err := t.variablesRepository.GetVariables(ctx)
	if err != nil {
		hlog.Error(ctx, "templateMediator.Save", err.Error())
		return nil, err
	}

	var templateVariables []string
	// for each component, header, body and footer
	re := regexp.MustCompile(`\{\{([^}]+)\}\}`)
	for i := range templateDomain.Components {
		v := templateDomain.Components[i]
		matches := re.FindAllStringSubmatch(*v.Text, -1)

		if matches != nil {
			frontVariable := matches[0][1]

			if _, ok := variables[frontVariable]; !ok {
				hlog.Error(ctx, "templateMediator.Save", "variable not found")
				return nil, fmt.Errorf("variable %s not found", frontVariable)
			}
			templateVariables = append(templateVariables, frontVariable)
			manyExamples := v.Example[strings.ToLower(fmt.Sprintf("%s_text", v.Type))]

			for j, eachExample := range manyExamples {
				if len(eachExample) != len(matches) {
					return nil, errors.New("error when save template")
				}
				replacedString := strings.Replace(*v.Text, frontVariable, strconv.Itoa(j+1), -1)
				templateDomain.Components[i].Text = &replacedString
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
	templateDomain.Components = cloned
	templateDomain.Variables = templateVariables
	saveTemplate, err := t.templateRepository.SaveTemplate(ctx, user, templateDomain)
	if err != nil {
		hlog.Error(ctx, "templateMediator.Save", fmt.Sprintf("%s", err))
		return nil, err
	}
	return saveTemplate, nil
}
