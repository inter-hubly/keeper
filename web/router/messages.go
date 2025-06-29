package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/infraestructure/middleware"
	"github.com/inter-hubly/keeper/internal/messages"
)

func newMessagesRouter(ctx context.Context, e *gin.RouterGroup) {
	controller := messages.NewController(ctx)
	messageGroup := e.Group("/messages").Use(middleware.AuthMiddleware())
	messageGroup.GET("/search", controller.SearchMessages)

}
