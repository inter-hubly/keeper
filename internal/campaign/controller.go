package campaign

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/infraestructure/middleware"
	"github.com/inter-hubly/keeper/internal/domain/kdto"
	"github.com/inter-hubly/keeper/web/httprest"
)

type Controller interface {
	GetCampaign(c *gin.Context)
	SaveCampaign(c *gin.Context)
	StartCampaign(c *gin.Context)
	SearchCampaigns(c *gin.Context)
}

type campaignController struct {
	campaignService Service
}

var (
	_campaignControllerOnce sync.Once
	_campaignController     *campaignController
)

func NewController(ctx context.Context) *campaignController {
	_campaignControllerOnce.Do(func() {
		_campaignController = &campaignController{
			campaignService: NewService(ctx),
		}
	})
	return _campaignController
}

func (ctrl *campaignController) GetCampaign(c *gin.Context) {
	ctx, loggedUser := middleware.GetLoggedUser(c)

	getCampaign, err := ctrl.campaignService.GetCampaign(ctx, loggedUser)

	if err != nil {
		httprest.Error(c, err.Error())
		return
	}

	httprest.Created(c, getCampaign)
}

func (ctrl *campaignController) StartCampaign(c *gin.Context) {
	ctx, loggedUser := middleware.GetLoggedUser(c)
	campaignId := c.Param("campaignId")

	if err := ctrl.campaignService.StartCampaign(ctx, loggedUser, campaignId); err != nil {
		httprest.Error(c, err.Error())
		return
	}

	httprest.Ok(c, nil)
}

func (ctrl *campaignController) SaveCampaign(c *gin.Context) {
	ctx, loggedUser := middleware.GetLoggedUser(c)

	var campaignDto kdto.Campaign

	if err := c.BindJSON(&campaignDto); err != nil {
		httprest.Error(c, "Error when marshal campaign")
		return
	}
	saveCampaign, err := ctrl.campaignService.SaveCampaign(ctx, loggedUser, &campaignDto)

	if err != nil {
		httprest.Error(c, err.Error())
		return
	}

	httprest.Created(c, saveCampaign)
}

func (ctrl *campaignController) SearchCampaigns(c *gin.Context) {
	ctx, loggedUser := middleware.GetLoggedUser(c)

	listCampaign, err := ctrl.campaignService.ListCampaign(ctx, loggedUser)
	if err != nil {
		httprest.Error(c, err.Error())
		return
	}
	httprest.Ok(c, listCampaign)
}
