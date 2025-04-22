//go:build e2e

package mediator

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/inter-hubly/keeper/internal/app/domain"
	"github.com/inter-hubly/keeper/internal/app/domain/kdto"
	gatewayMock "github.com/inter-hubly/keeper/internal/app/gateway/mocks"
	"github.com/inter-hubly/keeper/internal/app/repository/mocks"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/testutils"
	"github.com/stretchr/testify/assert"
)

type allMock struct {
	templateRepository  *mocks.MockTemplate
	variablesRepository *mocks.MockVariable
	whatsAppGateway     *gatewayMock.MockWhatsApp
}

func TestTemplateMediator(t *testing.T) {
	ctx := testutils.SetLoggedUser(context.Background())
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	allMocks := allMock{
		variablesRepository: mocks.NewMockVariable(ctrl),
		templateRepository:  mocks.NewMockTemplate(ctrl),
		whatsAppGateway:     gatewayMock.NewMockWhatsApp(ctrl),
	}

	mediator := templateMediator{
		variablesRepository: allMocks.variablesRepository,
		templateRepository:  allMocks.templateRepository,
		whatsAppGateway:     allMocks.whatsAppGateway,
	}

	loggedUser := hctx.LoggedUser.Get(ctx)
	variable := newVariables()

	for _, v := range []struct {
		testName string
		hasError bool
		template *domain.Template
		auxFunc  func(template *domain.Template)
	}{
		{
			testName: "Need to save new template",
			template: newDomainTemplate("Oi, {{name}}! Tudo bem? 😊 Só passando para lembrar sobre o pagamento da aula de inglês"),
			auxFunc: func(template *domain.Template) {
				allMocks.whatsAppGateway.EXPECT().CreateTemplate(gomock.Any(), gomock.Any()).
					Do(func(_ context.Context, template *domain.Template) {
						bodyText := "Oi, {{1}}! Tudo bem? 😊 Só passando para lembrar sobre o pagamento da aula de inglês"
						assert.Equal(t, template.Components[1].Text, &bodyText)
					}).
					Return(&kdto.WhatsAppTemplateResponse{
						Id: "123456",
					}, nil)

				allMocks.variablesRepository.EXPECT().GetVariables(gomock.Any()).Return(variable.Variable, nil)
				allMocks.templateRepository.EXPECT().SaveTemplate(gomock.Any(), gomock.Any(), gomock.Any()).Return(template, nil)
			},
		},
		{
			testName: "Cant save new template because does not exist variable",
			hasError: true,
			template: newDomainTemplate("Oi, {{testError}}! Tudo bem? 😊 Só passando para lembrar sobre o pagamento da aula de inglês"),
			auxFunc: func(template *domain.Template) {
				allMocks.variablesRepository.EXPECT().GetVariables(gomock.Any()).Return(variable.Variable, nil)
			},
		},
	} {
		t.Run(v.testName, func(t *testing.T) {
			v.auxFunc(v.template)
			save, err := mediator.Save(ctx, &loggedUser, v.template)
			if v.hasError {
				assert.Error(t, err)
				return
			}
			assert.NotEmpty(t, save)
			assert.Nil(t, err)
		})
	}
}

func newDomainTemplate(bodyText string) *domain.Template {
	headerText := "Lembrete de Pagamento"
	footerText := "See ya!"
	components := []domain.Components{
		{Type: "HEADER", Format: "TEXT", Text: &headerText},
		{Type: "BODY", Format: "TEXT", Text: &bodyText, Example: map[string][][]string{
			"body_text": {
				{"Saimon"},
			},
		},
		},
		{Type: "FOOTER", Format: "TEXT", Text: &footerText},
	}

	return &domain.Template{
		Name:            "test",
		Language:        "pt_BR",
		Category:        "UTILITY",
		ParameterFormat: "POSITIONAL",
		Components:      components,
	}
}

func newVariables() *domain.Variables {
	return &domain.Variables{
		Variable: map[string]domain.SingleVariable{
			"name": {
				Slug:  "name",
				Label: "name",
			},
		},
	}
}
