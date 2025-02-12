package controller

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/internal/app/domain"
	"github.com/inter-hubly/keeper/internal/app/domain/kdto"
	"github.com/inter-hubly/keeper/internal/app/service"
	"github.com/inter-hubly/keeper/internal/infraestructure/httprest"
	"github.com/inter-hubly/keeper/internal/infraestructure/middleware"
)

type Contact interface {
	SaveContact(c *gin.Context)
	FindContacts(c *gin.Context)
}

type contactController struct {
	contactService service.Contact
}

var (
	contactOnce sync.Once
	contact     *contactController
)

func NewContact(ctx context.Context) *contactController {
	contactOnce.Do(func() {
		contact = &contactController{
			contactService: service.NewContact(ctx),
		}
	})
	return contact
}

func (cc *contactController) SaveContact(c *gin.Context) {
	ctx, loggedUser := middleware.GetLoggedUser(c)

	var contactDto kdto.Contact

	if err := c.BindJSON(&contactDto); err != nil {
		httprest.Error(c, "Error when marshal body")
		return
	}

	var resp *domain.Contact
	var err error

	if resp, err = cc.contactService.SaveContact(ctx, loggedUser, &contactDto); err != nil {
		httprest.Error(c, "Error when save variable")
		return
	}
	httprest.Created(c, resp)
}

func (cc *contactController) FindContacts(c *gin.Context) {
	ctx, loggedUser := middleware.GetLoggedUser(c)

	resp, err := cc.contactService.FindContacts(ctx, loggedUser)

	if err != nil {
		httprest.Error(c, "Error when save variable")
		return
	}

	httprest.Ok(c, resp)
}
