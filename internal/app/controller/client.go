package controller

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/internal/app/service"
	"github.com/inter-hubly/keeper/internal/infraestructure/httprest"
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
	id := c.Param("id")
	getClient, err := ctrl.clientService.GetClientByPhoneNumberId(c, id)
	if err != nil {
		httprest.Error(c, "client not found")
		return
	}
	httprest.Ok(c, getClient)
}
