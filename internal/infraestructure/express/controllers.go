package express

import (
	"context"
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
	engine              *gin.Engine
	clientController    controller.Client
	messageController   controller.Messages
	authController      controller.Auth
	campaignController  controller.Campaign
	variablesController controller.Variables
	contactController   controller.Contact
	templateController  controller.Templates
}

func NewKeeperController(ctx context.Context, engine *gin.Engine) {
	controllersOnce.Do(func() {
		keeperControllers = &controllers{
			engine:              engine,
			clientController:    controller.NewClient(ctx),
			messageController:   controller.NewMessages(ctx),
			authController:      controller.NewAuth(ctx),
			campaignController:  controller.NewCampaign(ctx),
			variablesController: controller.NewVariable(ctx),
			contactController:   controller.NewContact(ctx),
			templateController:  controller.NewTemplate(ctx),
		}
	})
	keeperControllers.startControllers()
}

func (c *controllers) startControllers() {
	apiGroup := c.engine.Group("/api")
	{
		apiGroup.POST("/login", c.authController.Login)
		// apiGroup.POST("/sign-up", c.authController.CreateUser)
	}
	{

		clientGroup := apiGroup.Group("/client").Use(middleware.AuthMiddleware())
		clientGroup.POST("", c.clientController.SaveClient)
		clientGroup.GET("/phone-number-id", c.clientController.GetClientPhoneNumberId)
	}
	{
		messageGroup := apiGroup.Group("/messages").Use(middleware.AuthMiddleware())
		messageGroup.GET("/search", c.messageController.SearchMessages)
	}
	{
		campaignGroup := apiGroup.Group("/campaign").Use(middleware.AuthMiddleware())
		campaignGroup.POST("/:campaignId/start", c.campaignController.StartCampaign)
		campaignGroup.GET("", c.campaignController.GetCampaign)
		campaignGroup.POST("", c.campaignController.SaveCampaign)
	}
	{
		variablesGroup := apiGroup.Group("/variables").Use(middleware.AuthMiddleware())
		variablesGroup.GET("", c.variablesController.GetVariables)
		variablesGroup.POST("", c.variablesController.SaveManyVariable)
	}
	{
		contactGroup := apiGroup.Group("/contact").Use(middleware.AuthMiddleware())
		contactGroup.GET("", c.contactController.FindContacts)
		contactGroup.POST("", c.contactController.SaveContact)
	}
	{
		templateGroup := apiGroup.Group("/template").Use(middleware.AuthMiddleware())
		templateGroup.POST("", c.templateController.Save)
		templateGroup.GET("/search", c.templateController.FindAll)
	}
}
