package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/internal/infraestructure/express"
	"github.com/inter-hubly/pilot/hlog"
	"github.com/inter-hubly/pilot/server"
)

func main() {
	server.FillConfigEnvironment()

	router := gin.Default()

	express.Start(router)

	hlog.Info("main", fmt.Sprintf("Server start in port %d", server.GetEnvironment().Port))
	if err := router.Run(fmt.Sprintf(":%d", server.GetEnvironment().Port)); err != nil {
		hlog.Error("main", "Failed to start server: %v", err)
	}
}
