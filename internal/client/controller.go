package client

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/infraestructure/middleware"
	"github.com/inter-hubly/keeper/internal/domain/kdto"
	"github.com/inter-hubly/keeper/web/httprest"
)

type Client interface {
	GetClientPhoneNumberId(c *gin.Context)
	GetClient(c *gin.Context)
	SaveClient(c *gin.Context)
}

type clientController struct {
	clientService Service
}

var (
	_clientControllerOnce sync.Once
	_clientController     *clientController
)

func NewController(ctx context.Context) *clientController {
	_clientControllerOnce.Do(func() {
		_clientController = &clientController{
			clientService: NewService(ctx),
		}
	})
	return _clientController
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
