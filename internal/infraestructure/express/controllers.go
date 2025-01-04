package express

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/internal/app/controller"
)

var (
	controllersOnce   sync.Once
	keeperControllers *controllers
)

type controllers struct {
	engine            *gin.Engine
	clientController  controller.Client
	messageController controller.Messages
	authController    controller.Auth
}

func NewKeeperController(engine *gin.Engine) {
	controllersOnce.Do(func() {
		keeperControllers = &controllers{
			engine:            engine,
			clientController:  controller.NewClient(),
			messageController: controller.NewMessages(),
			authController:    controller.NewAuth(),
		}
	})
	keeperControllers.startControllers()
}

func (c *controllers) startControllers() {
	apiGroup := c.engine.Group("/api")
	{
		apiGroup.POST("/login", c.authController.Login)
	}
	{
		clientGroup := apiGroup.Group("/client")
		clientGroup.GET("/:id/phone-number-id", c.clientController.GetClientByPhoneNumberId)
	}
	{
		messageGroup := apiGroup.Group("/messages")
		messageGroup.POST("/search", c.messageController.SearchMessages)
	}
}
