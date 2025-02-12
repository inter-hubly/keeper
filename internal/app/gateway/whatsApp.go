package gateway

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/inter-hubly/keeper/internal/app/domain"
	"github.com/inter-hubly/keeper/internal/app/repository"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
	"github.com/inter-hubly/pilot/hrest"
)

type WhatsApp interface {
	CreateTemplate(ctx context.Context, message *domain.Template)
}

type whatsAppGateway struct {
	url              string
	clientRepository repository.Client
}

var (
	whatsAppOnce sync.Once
	whatsApp     *whatsAppGateway
)

func NewWhatsApp(ctx context.Context) *whatsAppGateway {
	whatsAppOnce.Do(func() {
		whatsApp = &whatsAppGateway{
			url:              "https://graph.facebook.com/v21.0/",
			clientRepository: repository.NewClient(ctx),
		}
	})
	return whatsApp
}

func (w *whatsAppGateway) CreateTemplate(ctx context.Context, message *domain.Template) {
	tenantId := hctx.Tenant.Get(ctx)
	client, err := w.clientRepository.GetClientByPhoneNumberId(ctx, tenantId)
	if err != nil {
		hlog.Error(ctx, "whatsAppGateway.CreateTemplate", err.Error())
	}
	request := hrest.NewRequest(fmt.Sprintf("%s/%s/message_template", w.url, client.AppId))

	if err = request.CreateRequest(ctx, http.MethodPost); err != nil {
		hlog.Error(ctx, "whatsAppGateway.CreateTemplate", err.Error())
	}

}
