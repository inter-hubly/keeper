package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/internal/variables"
	"github.com/inter-hubly/keeper/web/middleware"
)

func newVariablesRouter(ctx context.Context, e *gin.RouterGroup) {
	controller := variables.NewController(ctx)

	variablesGroup := e.Group("/variables").Use(middleware.AuthMiddleware())
	variablesGroup.GET("/search", controller.SearchVariables)
	variablesGroup.POST("", controller.SaveManyVariable)

}
