package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/pilot/database/hmongo"
	"github.com/inter-hubly/pilot/domain/entity"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Campaign interface {
	GetCampaignById(ctx context.Context, campaignId string) (*entity.Campaign, error)
	SaveCampaign(ctx context.Context, campaign *entity.Campaign) (*entity.Campaign, error)
	ListCampaign(ctx context.Context) ([]entity.Campaign, error)
}

type campaignRepository struct {
	connection hmongo.NoSqlConn
	collection string
}

var (
	onceCampaign       sync.Once
	repositoryCampaign *campaignRepository
)

func NewCampaign(ctx context.Context) *campaignRepository {
	onceCampaign.Do(func() {
		repositoryCampaign = &campaignRepository{
			connection: hmongo.GetConnection(ctx),
			collection: "campaign",
		}
	})
	return repositoryCampaign
}

func (c *campaignRepository) GetCampaignById(ctx context.Context, campaignId string) (*entity.Campaign, error) {
	hlog.Debug(ctx, "campaignRepository.GetCampaignById", fmt.Sprintf("campaignId: %s", campaignId))
	var campaign entity.Campaign

	tenantId := hctx.Tenant.Get(ctx)

	if err := c.connection.GetCollection(ctx, c.collection).FindOne(ctx,
		bson.M{
			"tenantId": tenantId,
			"_id":      campaignId,
		},
	).Decode(&campaign); err != nil {
		hlog.Error(ctx, "campaignRepository.GetCampaignById", fmt.Sprintf("campaignId: %s", campaignId))
		return nil, err
	}
	return &campaign, nil
}

func (c *campaignRepository) SaveCampaign(ctx context.Context, campaign *entity.Campaign) (*entity.Campaign, error) {
	hlog.Debug(ctx, "campaignRepository.SaveCampaign", fmt.Sprintf("saving campaign %s", campaign.Id))
	one, err := c.connection.GetCollection(ctx, c.collection).InsertOne(ctx, campaign)
	if err != nil {
		hlog.Error(ctx, "campaignRepository.SaveCampaign", fmt.Sprintf("error while saving campaign %s", campaign.Id))
		return nil, err
	}
	id := one.InsertedID.(primitive.ObjectID)
	campaign.Id = id.Hex()
	return campaign, nil
}

func (c *campaignRepository) ListCampaign(ctx context.Context) ([]entity.Campaign, error) {
	hlog.Debug(ctx, "campaignRepository.ListCampaign", fmt.Sprintf("list campaign"))
	tenantId := hctx.Tenant.Get(ctx)

	cur, err := c.connection.GetCollection(ctx, c.collection).Find(ctx, bson.M{
		"tenantId": tenantId,
	})
	if err != nil {
		hlog.Error(ctx, "campaignRepository.ListCampaign", fmt.Sprintf("error while list campaign %s", tenantId))
		return nil, err
	}
	var campaigns []entity.Campaign
	if err = cur.All(ctx, &campaigns); err != nil {
		hlog.Error(ctx, "campaignRepository.ListCampaign", fmt.Sprintf("error while list campaign %s", tenantId))
	}
	return campaigns, nil
}
