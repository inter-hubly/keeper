package gateway

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/inter-hubly/keeper/internal/app/domain"
	"github.com/inter-hubly/keeper/internal/app/domain/kdto"
	"github.com/inter-hubly/keeper/internal/app/repository"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
	"github.com/inter-hubly/pilot/hrest"
	"github.com/inter-hubly/pilot/util"
)

type WhatsApp interface {
	CreateTemplate(ctx context.Context, message *domain.Template) (*kdto.WhatsAppTemplateResponse, error)
}

type whatsAppGateway struct {
	url              string
	clientRepository repository.Client
}

var (
	_whatsAppOnce sync.Once
	_whatsApp     *whatsAppGateway
)

func NewWhatsApp(ctx context.Context) *whatsAppGateway {
	_whatsAppOnce.Do(func() {
		_whatsApp = &whatsAppGateway{
			url:              "https://graph.facebook.com/v21.0",
			clientRepository: repository.NewClient(ctx),
		}
	})
	return _whatsApp
}

func (w *whatsAppGateway) CreateTemplate(ctx context.Context, message *domain.Template) (*kdto.WhatsAppTemplateResponse, error) {
	tenantId := hctx.Tenant.Get(ctx)
	client, err := w.clientRepository.GetClientByPhoneNumberId(ctx, tenantId)
	if err != nil {
		hlog.Error(ctx, "whatsAppGateway.CreateTemplate", err.Error())
	}

	messageDto := struct {
		Name            string         `json:"name"`
		Language        string         `json:"language"`
		Category        string         `json:"category"`
		ParameterFormat string         `json:"parameter_format"`
		Components      []componentDto `json:"components"`
	}{
		Name:            util.ToSlug(message.Name, true),
		Language:        message.Language,
		Category:        message.Category,
		ParameterFormat: message.ParameterFormat,
	}

	for _, component := range message.Components {
		dto := componentDto{
			Type: string(component.Type),
			Text: component.Text,
		}
		if len(component.Example) > 0 {
			dto.Example = component.Example
		}
		if component.Type == "HEADER" {
			dto.Format = "TEXT"
		}
		messageDto.Components = append(messageDto.Components, dto)
	}

	request := hrest.NewRequest(fmt.Sprintf("%s/%s/message_templates", w.url, client.BusinessId), hrest.WithBody(messageDto),
		hrest.WithHeader([]hrest.Pair[string, string]{
			{"Content-Type", "application/json"},
			{"Authorization", "Bearer " + client.AccessToken},
		}))

	if err = request.CreateRequest(ctx, http.MethodPost); err != nil {
		hlog.Error(ctx, "whatsAppGateway.CreateTemplate", err.Error())
		return nil, err
	}
	var whatsResp kdto.WhatsAppTemplateResponse
	if err = request.GetBody(ctx, &whatsResp); err != nil {
		hlog.Error(ctx, "whatsAppGateway.CreateTemplate", err.Error())
		return nil, err
	}
	return &whatsResp, nil
}

type componentDto struct {
	Type    string                `json:"type"`
	Text    string                `json:"text"`
	Format  string                `json:"format,omitempty"`
	Example map[string][][]string `json:"example,omitempty"`
}
