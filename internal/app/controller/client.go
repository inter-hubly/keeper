package controller

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/internal/app/domain/kdto"
	"github.com/inter-hubly/keeper/internal/app/service"
	"github.com/inter-hubly/keeper/internal/infraestructure/httprest"
	"github.com/inter-hubly/keeper/internal/infraestructure/middleware"
)

type Client interface {
	GetClientPhoneNumberId(c *gin.Context)
	GetClient(c *gin.Context)
	SaveClient(c *gin.Context)
}

type clientController struct {
	clientService service.Client
}

var (
	clientControllerOnce sync.Once
	client               *clientController
)

func NewClient(ctx context.Context) *clientController {
	clientControllerOnce.Do(func() {
		client = &clientController{
			clientService: service.NewClient(ctx),
		}
	})
	return client
}

func (ctrl *clientController) GetClientPhoneNumberId(c *gin.Context) {
	ctx, loggedUser := middleware.GetLoggedUser(c)
	getClient, err := ctrl.clientService.GetClientByPhoneNumberId(ctx, loggedUser.Tenant)
	if err != nil {
		httprest.Error(c, "client not found")
		return
	}
	httprest.Ok(c, struct {
		PhoneNumberId string `json:"phoneNumberId"`
	}{
		PhoneNumberId: getClient.PhoneNumberId,
	})
}

func (ctrl *clientController) SaveClient(c *gin.Context) {
	var clientDto kdto.Client

	if err := c.BindJSON(&clientDto); err != nil {
		httprest.Error(c, "Error when marshal body")
		return
	}

	ctx, loggedUser := middleware.GetLoggedUser(c)
	if err := ctrl.clientService.SaveClient(ctx, loggedUser, clientDto); err != nil {
		httprest.Error(c, err.Error())
		return
	}
	httprest.Created(c, nil)
}

func (ctrl *clientController) GetClient(c *gin.Context) {
	ctx, loggedUser := middleware.GetLoggedUser(c)
	getClient, err := ctrl.clientService.GetClientByPhoneNumberId(ctx, loggedUser.Tenant)
	if err != nil {
		httprest.Error(c, "client not found")
		return
	}
	httprest.Ok(c, getClient)
}
