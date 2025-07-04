//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go -package=mocks

package campaign

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/inter-hubly/pilot/database/hmongo"
	"github.com/inter-hubly/pilot/domain/entity"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	GetCampaignById(ctx context.Context, campaignId string) (*entity.Campaign, error)
	GetCampaignsByIds(ctx context.Context, campaignId ...string) ([]entity.Campaign, error)
	GetCampaignsTemplateById(ctx context.Context, campaignId []string) ([]string, error)
	SaveCampaign(ctx context.Context, campaign *entity.Campaign) error
	ListCampaign(ctx context.Context) ([]entity.Campaign, error)
}

type campaignRepository struct {
	connection hmongo.NoSqlConn
	collection string
}

var (
	_campaignRepositoryOnce sync.Once
	_campaignRepository     *campaignRepository
)

func NewRepository(ctx context.Context) *campaignRepository {
	_campaignRepositoryOnce.Do(func() {
		_campaignRepository = &campaignRepository{
			connection: hmongo.GetConnection(ctx),
			collection: "campaigns",
		}
	})
	return _campaignRepository
}

func (c *campaignRepository) GetCampaignById(ctx context.Context, campaignId string) (*entity.Campaign, error) {
	hlog.Debug(ctx, "campaign.repository.GetCampaignById", fmt.Sprintf("campaignId: %s", campaignId))
	var campaign entity.Campaign

	tenantId := hctx.Tenant.Get(ctx)

	if err := c.connection.GetCollection(ctx, c.collection).FindOne(ctx,
		bson.M{
			"tenantId": tenantId,
			"_id":      campaignId,
		},
	).Decode(&campaign); err != nil {
		hlog.Error(ctx, "campaign.repository.GetCampaignById", fmt.Sprintf("campaignId: %s", campaignId))
		return nil, err
	}
	return &campaign, nil
}

func (c *campaignRepository) SaveCampaign(ctx context.Context, campaign *entity.Campaign) error {
	if campaign == nil {
		return errors.New("campaign is nil")
	}

	hlog.Debug(ctx, "campaign.repository.SaveCampaign", fmt.Sprintf("saving campaign %s", campaign.Id))
	one, err := c.connection.GetCollection(ctx, c.collection).InsertOne(ctx, campaign)
	if err != nil {
		hlog.Error(ctx, "campaign.repository.SaveCampaign", fmt.Sprintf("error while saving campaign %s", campaign.Id))
		return err
	}
	id := one.InsertedID.(primitive.ObjectID)
	campaign.Id = id.Hex()
	return nil
}

func (c *campaignRepository) ListCampaign(ctx context.Context) ([]entity.Campaign, error) {
	hlog.Debug(ctx, "campaign.repository.ListCampaign", fmt.Sprintf("list campaign"))
	tenantId := hctx.Tenant.Get(ctx)

	cur, err := c.connection.GetCollection(ctx, c.collection).Find(ctx, bson.M{
		"tenantId": tenantId,
	})
	if err != nil {
		hlog.Error(ctx, "campaign.repository.ListCampaign", fmt.Sprintf("error while list campaign %s", tenantId))
		return nil, err
	}
	var campaigns []entity.Campaign
	if err = cur.All(ctx, &campaigns); err != nil {
		hlog.Error(ctx, "campaign.repository.ListCampaign", fmt.Sprintf("error while list campaign %s", tenantId))
	}
	return campaigns, nil
}

func (c *campaignRepository) GetCampaignsTemplateById(ctx context.Context, campaignId []string) ([]string, error) {
	hlog.Debug(ctx, "campaign.repository.GetCampaignsById", fmt.Sprintf("campaignId: %s", campaignId))
	tenantId := hctx.Tenant.Get(ctx)
	var (
		campaign []entity.Campaign
		err      error
	)

	//var objectIDs []primitive.ObjectID
	//for _, id := range campaignId {
	//	objID, err := primitive.ObjectIDFromHex(id)
	//	if err != nil {
	//		hlog.Error(ctx, "campaign.repository.GetCampaignsById", fmt.Sprintf("invalid ObjectId: %s", id))
	//		return nil, err
	//	}
	//	objectIDs = append(objectIDs, objID)
	//}

	cur, err := c.connection.GetCollection(ctx, c.collection).Find(ctx, bson.M{
		"_id": bson.M{"$in": campaignId},
	}, options.Find().SetProjection(
		bson.M{
			"template.id": 1,
			"_id":         0,
		}))

	if err != nil {
		hlog.Error(ctx, "campaign.repository.GetCampaign sById", fmt.Sprintf("error while list campaign %s", tenantId))
		return nil, err
	}
	if err = cur.All(ctx, &campaign); err != nil {
		hlog.Error(ctx, "campaign.repository.GetCampaignsById", fmt.Sprintf("error while list campaign %s", tenantId))
		return nil, err
	}
	var response []string
	for _, objID := range campaign {
		response = append(response, objID.Template.Id)
	}
	return response, nil
}

func (c *campaignRepository) GetCampaignsByIds(ctx context.Context, campaignId ...string) ([]entity.Campaign, error) {
	hlog.Debug(ctx, "campaign.repository.GetCampaignsByIds", fmt.Sprintf("campaignId: %s", campaignId))
	var campaigns []entity.Campaign
	//objectIDs []primitive.ObjectID

	tenantId := hctx.Tenant.Get(ctx)
	//for _, id := range campaignId {
	//	objID, err := primitive.ObjectIDFromHex(id)
	//	if err != nil {
	//		hlog.Error(ctx, "campaign.repository.GetCampaignsById", fmt.Sprintf("invalid ObjectId: %s", id))
	//		return nil, err
	//	}
	//	objectIDs = append(objectIDs, objID)
	//}
	cur, err := c.connection.GetCollection(ctx, c.collection).Find(ctx, bson.M{
		"_id":      bson.M{"$in": campaignId},
		"tenantId": tenantId,
	})
	if err != nil {
		hlog.Error(ctx, "campaign.repository.GetCampaignsById", fmt.Sprintf("error while list campaign %s", tenantId))
		return nil, err
	}
	if err = cur.All(ctx, &campaigns); err != nil {
		hlog.Error(ctx, "campaign.repository.GetCampaignsById", fmt.Sprintf("error while list campaign %s", tenantId))
		return nil, err
	}
	return campaigns, nil
}
