package express

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/internal/app/controller"
	"github.com/inter-hubly/keeper/internal/infraestructure/middleware"
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
		apiGroup.POST("/sign-up", c.authController.CreateUser)
	}
	{

		clientGroup := apiGroup.Group("/client").Use(middleware.AuthMiddleware())
		clientGroup.GET("/phone-number-id", c.clientController.GetClientPhoneNumberId)
	}
	{
		messageGroup := apiGroup.Group("/messages").Use(middleware.AuthMiddleware())
		messageGroup.POST("/search", c.messageController.SearchMessages)
	}
}
