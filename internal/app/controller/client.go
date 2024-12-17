package controller

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/internal/app/service"
)

type Client interface {
	GetClient(c *gin.Context)
}

var (
	clientControllerOnce sync.Once
	client               *clientController
)

type clientController struct {
	clientService service.Client
}

func NewClient() *clientController {
	clientControllerOnce.Do(func() {
		client = &clientController{}
	})
	return client
}

func (ctrl *clientController) GetClient(c *gin.Context) {
	id := c.Param("id")
	getClient, err := ctrl.clientService.GetClient(c, id)
	if err != nil {
		c.JSON(500, err)
	}
	c.JSON(200, getClient)
}
