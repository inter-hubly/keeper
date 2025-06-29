package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/internal/contact"
	"github.com/inter-hubly/keeper/web/middleware"
)

func newContactRouter(ctx context.Context, e *gin.RouterGroup) {
	controller := contact.NewController(ctx)

	contactGroup := e.Group("/contact").Use(middleware.AuthMiddleware())
	contactGroup.GET("", controller.SearchContacts)
	contactGroup.POST("", controller.SaveContact)
}
