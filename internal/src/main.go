package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/internal/infraestructure/express"
	"github.com/inter-hubly/pilot/hlog"
	"github.com/inter-hubly/pilot/server"
)

func main() {
	ctx := context.Background()
	server.FillConfigEnvironment(ctx)
	router := gin.Default()

	express.Start(ctx, router)

	hlog.Info(ctx, "main", fmt.Sprintf("Server start in port %s", server.GetEnvironment().Port))
	if err := router.Run(fmt.Sprintf(":%s", server.GetEnvironment().Port)); err != nil {
		hlog.Error(ctx, "main", fmt.Sprintf("Failed to start server: %v", err))
	}
}
