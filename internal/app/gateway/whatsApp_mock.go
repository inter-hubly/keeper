package gateway

import (
	"context"

	"github.com/inter-hubly/keeper/internal/app/domain"
	"github.com/inter-hubly/keeper/internal/app/domain/kdto"
)

type WhatsAppMock struct {
}

func NewWhatsAppMock() *WhatsAppMock {
	return &WhatsAppMock{}
}

func (w *WhatsAppMock) CreateTemplate(ctx context.Context, message *domain.Template) (*kdto.WhatsAppTemplateResponse, error) {
	return &kdto.WhatsAppTemplateResponse{
		Id:       "942970764616888",
		Status:   "APPROVED",
		Category: "UTILITY",
	}, nil
}
