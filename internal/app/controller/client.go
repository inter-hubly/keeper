package controller

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/internal/app/service"
	"github.com/inter-hubly/keeper/internal/infraestructure/httprest"
	"github.com/inter-hubly/keeper/internal/infraestructure/middleware"
)

type Client interface {
	GetClientByPhoneNumberId(c *gin.Context)
}

type clientController struct {
	clientService service.Client
}

func NewClient() *clientController {

	var (
		clientControllerOnce sync.Once
		client               *clientController
	)

	clientControllerOnce.Do(func() {
		client = &clientController{
			clientService: service.NewClient(),
		}
	})
	return client
}

func (ctrl *clientController) GetClientByPhoneNumberId(c *gin.Context) {
	ctx, loggedUser := middleware.GetLoggedUser(c)
	getClient, err := ctrl.clientService.GetClientByPhoneNumberId(ctx, loggedUser.Tenant)
	if err != nil {
		httprest.Error(c, "client not found")
		return
	}
	httprest.Ok(c, getClient)
}
