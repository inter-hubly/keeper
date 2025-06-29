package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/infraestructure/middleware"
	"github.com/inter-hubly/keeper/internal/contact"
)

func newContactRouter(ctx context.Context, e *gin.RouterGroup) {
	controller := contact.NewController(ctx)

	contactGroup := e.Group("/contact").Use(middleware.AuthMiddleware())
	contactGroup.GET("", controller.SearchContacts)
	contactGroup.POST("", controller.SaveContact)
}
