package mediator

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/inter-hubly/keeper/internal/app/domain"
	"github.com/inter-hubly/keeper/internal/app/gateway"
	"github.com/inter-hubly/keeper/internal/app/repository"
	"github.com/inter-hubly/pilot/database/hmongo"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/testutils"
	"github.com/stretchr/testify/assert"
)

const containerEnvironments = false

func TestTemplateMediator(t *testing.T) {
	ctx := context.Background()

	var host string
	var close func(ctx context.Context) error
	var err error

	if containerEnvironments {
		host, close, err = testutils.ElasticSearch(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if close != nil {
			defer close(ctx)
		}
	} else {
		host = "mongodb://localhost:27017"
	}
	hmongo.NewConnection(
		ctx,
		hmongo.WithDatabase("test"),
		hmongo.WithUrl(host),
	)
	logged := hctx.Logged{
		UserId: "userTest",
		Tenant: "tenantTest",
	}
	ctx = hctx.LoggedUser.New(logged)

	mediator := templateMediator{
		templateRepository: repository.NewTemplate(ctx),
		whatsAppGateway:    gateway.NewWhatsAppMock(),
	}

	t.Run("Need to save", func(t *testing.T) {
		var myDto domain.Template

		if err = json.Unmarshal([]byte(dataDto), &myDto); err != nil {
			fmt.Println("Erro ao decodificar JSON:", err)
			return
		}

		save, err := mediator.Save(ctx, &logged, &myDto)
		assert.NotEmpty(t, save)
		assert.Nil(t, err)
	})
}

var dataDto = `{
		"name": "cobranca_mensal_1",
		"language": "pt_BR",
		"category": "UTILITY",
		"parameter_format": "POSITIONAL",
		"components": [
			{
				"type": "HEADER",
				"format": "TEXT",
				"text": "Lembrete de Pagamento"
			},
			{
				"type": "BODY",
				"text": "Oi, {{1}}! Tudo bem? 😊 Só passando para lembrar sobre o pagamento da aula de inglês de R$ {{2}} 📝. Assim a gente mantém tudo certinho e posso continuar te ajudando a arrasar no inglês! 🚀 Qualquer dúvida, é só me chamar! 😄",
				"example": {
					"body_text": [
						["Saimon", "20,20"]
					]
				}
			},
			{
				"type": "FOOTER",
				"text": "See ya!"
			}
		]
}`
