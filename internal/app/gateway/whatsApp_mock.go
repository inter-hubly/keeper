package gateway

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/inter-hubly/keeper/internal/app/domain"
	"github.com/inter-hubly/keeper/internal/app/domain/kdto"
)

type WhatsAppMock struct {
}

func NewWhatsAppMock() *WhatsAppMock {
	return &WhatsAppMock{}
}

func (w *WhatsAppMock) CreateTemplate(ctx context.Context, message *domain.Template) (*kdto.WhatsAppTemplateResponse, error) {
	var tmp domain.Template
	if err := json.Unmarshal([]byte(resp), &tmp); err != nil {
		panic(err)
	}
	if tmp.Language != message.Language {
		return nil, errors.New("language not match")
	}
	if tmp.Name != message.Name {
		return nil, errors.New("name not match")
	}
	if tmp.Category != message.Category {
		return nil, errors.New("category not match")
	}
	if tmp.ParameterFormat != message.ParameterFormat {
		return nil, errors.New("parameter_format not match")
	}
	if len(tmp.Components) != len(message.Components) {
		return nil, errors.New("components count not match")
	}

	return &kdto.WhatsAppTemplateResponse{
		Id:       "942970764616888",
		Status:   "APPROVED",
		Category: "UTILITY",
	}, nil
}

const resp = `{
    "name": "cobranca_mensal_1",
    "category": "UTILITY",
    "parameterFormat": "POSITIONAL",
    "language": "pt_BR",
    "components": [
        {
            "type": "HEADER",
            "text": "Lembrete de Pagamento",
            "format": "TEXT"
        },
        {
            "type": "BODY",
            "text": "Oi, {{1}}! Tudo bem? 😊 Só passando para lembrar sobre o pagamento da aula de inglês de R$ {{2}} 📝. Assim a gente mantém tudo certinho e posso continuar te ajudando a arrasar no inglês! 🚀 Qualquer dúvida, é só me chamar! 😄",
            "format": "TEXT",
            "example": {
                "body_text": [
                    [
                        "Saimon",
                        "20.20"
                    ]
                ]
            }
        },
        {
            "type": "FOOTER",
            "text": "See ya!",
            "format": "TEXT"
        }
    ]
}`
