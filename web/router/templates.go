package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/infraestructure/middleware"
	"github.com/inter-hubly/keeper/internal/templates"
)

func newTemplatesRouter(ctx context.Context, e *gin.RouterGroup) {
	controller := templates.NewController(ctx)

	templateGroup := e.Group("/templates").Use(middleware.AuthMiddleware())
	templateGroup.POST("", controller.Save)
	templateGroup.GET("/search", controller.SearchTemplates)
	templateGroup.POST("/sincronize", controller.SincronizeWhatsAppTemplate)
	templateGroup.POST("/:templateId/variables", controller.SaveVariables)
	templateGroup.GET("/:templateId/variables/count", controller.CountVariables)
}
