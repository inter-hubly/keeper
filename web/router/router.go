package router

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

var apiBasicRout = "/api/v1"

func NewRouter(ctx context.Context) *gin.Engine {
	router := gin.Default()
	router.Use(corsMiddleware())
	group := router.Group(apiBasicRout)

	newUserRouter(ctx, group)
	newClientRouter(ctx, group)
	newCampaignRouter(ctx, group)
	newMessagesRouter(ctx, group)
	newVariablesRouter(ctx, group)
	newTemplatesRouter(ctx, group)
	newContactRouter(ctx, group)

	return router
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, tenant")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
