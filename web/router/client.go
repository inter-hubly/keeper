package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/infraestructure/middleware"
	"github.com/inter-hubly/keeper/internal/client"
)

func newClientRouter(ctx context.Context, e *gin.RouterGroup) {
	controller := client.NewController(ctx)

	group := e.Group("/clients").Use(middleware.AuthMiddleware())
	group.POST("", controller.SaveClient)
	group.GET("", controller.GetClient)
	group.GET("/phone-number-id", controller.GetClientPhoneNumberId)
}
