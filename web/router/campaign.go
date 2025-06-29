package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/internal/campaign"
	"github.com/inter-hubly/keeper/web/middleware"
)

func newCampaignRouter(ctx context.Context, e *gin.RouterGroup) {
	controller := campaign.NewController(ctx)

	campaignGroup := e.Group("/campaigns").Use(middleware.AuthMiddleware())
	campaignGroup.POST("/:campaignId/start", controller.StartCampaign)
	campaignGroup.GET("/search", controller.SearchCampaigns)
	campaignGroup.GET("", controller.GetCampaign)
	campaignGroup.POST("", controller.SaveCampaign)

}
