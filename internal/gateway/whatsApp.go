//go:generate mockgen -source=whatsApp.go -destination=mocks/whatsApp_mock.go -package=mocks

package gateway

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/inter-hubly/keeper/internal/client"
	"github.com/inter-hubly/keeper/internal/domain"
	"github.com/inter-hubly/keeper/internal/domain/kdto"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
	"github.com/inter-hubly/pilot/hrest"
)

type WhatsApp interface {
	CreateTemplate(ctx context.Context, message *domain.Template) (*kdto.WhatsAppTemplateResponse, error)
	FindAllTemplate(ctx context.Context) ([]kdto.WhatsAppMessageTemplateResponse, error)
}

type whatsAppGateway struct {
	url              string
	clientRepository client.Repository
}

var (
	_whatsAppOnce sync.Once
	_whatsApp     *whatsAppGateway
)

func NewWhatsApp(ctx context.Context) *whatsAppGateway {
	_whatsAppOnce.Do(func() {
		_whatsApp = &whatsAppGateway{
			url:              "https://graph.facebook.com/v21.0",
			clientRepository: client.NewRepository(ctx),
		}
	})
	return _whatsApp
}

func (w *whatsAppGateway) CreateTemplate(ctx context.Context, message *domain.Template) (*kdto.WhatsAppTemplateResponse, error) {
	tenantId := hctx.Tenant.Get(ctx)
	client, err := w.clientRepository.GetClientByPhoneNumberId(ctx, tenantId)
	if err != nil {
		hlog.Error(ctx, "whatsAppGateway.CreateTemplate", err.Error())
		return nil, err
	}

	messageDto := struct {
		Name            string         `json:"name"`
		Language        string         `json:"language"`
		Category        string         `json:"category"`
		ParameterFormat string         `json:"parameter_format"`
		Components      []componentDto `json:"components"`
	}{
		Name:            message.Slug,
		Language:        message.Language,
		Category:        message.Category,
		ParameterFormat: message.ParameterFormat,
	}

	for _, component := range message.Components {
		dto := componentDto{
			Type: string(component.Type),
			Text: *component.Text,
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

func (w *whatsAppGateway) FindAllTemplate(ctx context.Context) ([]kdto.WhatsAppMessageTemplateResponse, error) {
	hlog.Debug(ctx, "whatsAppGateway.FindAllTemplate", "Find all template")

	tenantId := hctx.Tenant.Get(ctx)
	client, err := w.clientRepository.GetClientByPhoneNumberId(ctx, tenantId)
	if err != nil {
		hlog.Error(ctx, "whatsAppGateway.CreateTemplate", err.Error())
		return nil, err
	}

	request := hrest.NewRequest(fmt.Sprintf("%s/%s/message_templates", w.url, client.BusinessId),
		hrest.WithHeader([]hrest.Pair[string, string]{
			{"Content-Type", "application/json"},
			{"Authorization", "Bearer " + client.AccessToken},
		}))

	if err = request.CreateRequest(ctx, http.MethodGet); err != nil {
		hlog.Error(ctx, "whatsAppGateway.CreateTemplate", err.Error())
		return nil, err
	}
	var whatsResp struct {
		Data []kdto.WhatsAppMessageTemplateResponse `json:"data"`
	}
	if err = request.GetBody(ctx, &whatsResp); err != nil {
		hlog.Error(ctx, "whatsAppGateway.CreateTemplate", err.Error())
		return nil, err
	}
	return whatsResp.Data, nil
}

type componentDto struct {
	Type    string                `json:"type"`
	Text    string                `json:"text"`
	Format  string                `json:"format,omitempty"`
	Example map[string][][]string `json:"example,omitempty"`
}
