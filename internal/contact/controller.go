package contact

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/infraestructure/middleware"
	"github.com/inter-hubly/keeper/internal/domain"
	"github.com/inter-hubly/keeper/internal/domain/kdto"
	"github.com/inter-hubly/keeper/web/httprest"
)

type Controller interface {
	SaveContact(c *gin.Context)
	SearchContacts(c *gin.Context)
}

type contactController struct {
	contactService Service
}

var (
	_contactControllerOnce sync.Once
	_contactController     *contactController
)

func NewController(ctx context.Context) *contactController {
	_contactControllerOnce.Do(func() {
		_contactController = &contactController{
			contactService: NewService(ctx),
		}
	})
	return _contactController
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

func (cc *contactController) SearchContacts(c *gin.Context) {
	ctx, loggedUser := middleware.GetLoggedUser(c)

	resp, err := cc.contactService.SearchContacts(ctx, loggedUser)

	if err != nil {
		httprest.Error(c, "Error when save variable")
		return
	}

	httprest.Ok(c, resp)
}
