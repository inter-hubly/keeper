package contact

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/inter-hubly/keeper/internal/domain"
	"github.com/inter-hubly/keeper/internal/domain/kdto"
	"github.com/inter-hubly/keeper/internal/variables"
	"github.com/inter-hubly/pilot/domain/base"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
)

type Service interface {
	SaveContact(ctx context.Context, loggedUser *hctx.Logged, contactDto *kdto.Contact) (*domain.Contact, error)
	SearchContacts(ctx context.Context, loggedUser *hctx.Logged) ([]domain.Contact, error)
}

type contactService struct {
	contactRepository   Repository
	variablesRepository variables.Repository
}

var (
	_contactServiceOnce sync.Once
	_contactService     *contactService
)

func NewService(ctx context.Context) *contactService {
	_contactServiceOnce.Do(func() {
		_contactService = &contactService{
			contactRepository:   NewContact(ctx),
			variablesRepository: variables.NewRepository(ctx),
		}
	})
	return _contactService
}

func (c *contactService) SaveContact(ctx context.Context, loggedUser *hctx.Logged, contactDto *kdto.Contact) (*domain.Contact, error) {
	hlog.Debug(ctx, "contactService.SaveContact", fmt.Sprintf("%+v", contactDto))
	variableDb, err := c.variablesRepository.GetVariables(ctx)
	if err != nil {
		hlog.Error(ctx, "contactService.SaveContact", err.Error())
	}

	// verificando variável
	for _, dtoVariable := range contactDto.Variables {
		found := false
		for _, dbVariable := range variableDb {
			if dtoVariable.Key == dbVariable.Slug {
				found = true
			}
		}
		if !found {
			return nil, errors.New("variable not found")
		}
	}

	domainContact := &domain.Contact{
		Name:      contactDto.Name,
		Phone:     contactDto.Phone,
		Variables: contactDto.Variables,
		Entity:    base.NewBaseEntity(ctx, loggedUser),
	}

	contactId, err := c.contactRepository.SaveContact(ctx, domainContact)
	if err != nil {
		hlog.Error(ctx, "contactService.SaveContact", err.Error())
		return nil, err
	}
	domainContact.Id = contactId
	return domainContact, nil
}

func (c *contactService) SearchContacts(ctx context.Context, loggedUser *hctx.Logged) ([]domain.Contact, error) {
	hlog.Debug(ctx, "contactService.FindContacts", fmt.Sprintf("%+v", loggedUser))
	contacts, err := c.contactRepository.FindContacts(ctx, loggedUser.Tenant)
	if err != nil {
		return nil, err
	}
	return contacts, nil
}
