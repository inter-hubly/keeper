package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/internal/user"
)

func newUserRouter(ctx context.Context, e *gin.RouterGroup) {
	controller := user.NewController(ctx)

	group := e.Group("/users")
	group.POST("/login", controller.Login)
}
