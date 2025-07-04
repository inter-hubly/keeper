package campaign_test

import (
	"context"
	"errors"
	"github.com/inter-hubly/keeper/internal/campaign"
	"github.com/inter-hubly/pilot/domain/base"
	"testing"

	"github.com/inter-hubly/keeper/ktest"
	"github.com/inter-hubly/pilot/database/hmongo"
	"github.com/inter-hubly/pilot/domain/entity"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/testutils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRepository(t *testing.T) {
	ctx := testutils.SetLoggedUser(context.Background())

	mongoClose := ktest.MongoSetup(ctx)
	defer mongoClose(ctx)

	repository := campaign.NewRepository(ctx)
	insertedValues := setupDb(ctx)

	for _, tt := range []struct {
		name    string
		isError bool
		auxFunc func() error
	}{
		{
			name: "GetCampaignById need return success",
			auxFunc: func() error {
				campaignEntity := insertedValues[0]
				resp, err := repository.GetCampaignById(ctx, campaignEntity.Id)
				assert.Nil(t, err)
				assert.Equal(t, campaignEntity.Id, resp.Id)
				assert.Equal(t, campaignEntity.Name, resp.Name)
				return nil
			},
		},
		{
			name:    "GetCampaignById need return error",
			isError: true,
			auxFunc: func() error {
				resp, err := repository.GetCampaignById(ctx, "")
				assert.NotNil(t, err)
				assert.Nil(t, resp)
				return err
			},
		},
		{
			name:    "GetCampaignById need return error without tenant",
			isError: true,
			auxFunc: func() error {
				campaignEntity := insertedValues[0]
				resp, err := repository.GetCampaignById(context.Background(), campaignEntity.Id)
				assert.NotNil(t, err)
				assert.Nil(t, resp)
				return err
			},
		},
		{
			name: "GetCampaignsByIds need return success",
			auxFunc: func() error {
				ids := make([]string, 0, len(insertedValues))
				for _, id := range insertedValues {
					ids = append(ids, id.Id)
				}
				resp, err := repository.GetCampaignsByIds(ctx, ids...)
				assert.NotNil(t, resp)
				assert.Nil(t, err)
				assert.Equal(t, len(ids), len(insertedValues))
				return nil
			},
		},
		{
			name: "GetCampaignsByIds need return error because not found",
			auxFunc: func() error {
				ids := make([]string, 0, len(insertedValues))
				for _, id := range insertedValues {
					ids = append(ids, id.Id)
				}
				resp, err := repository.GetCampaignsByIds(context.Background(), ids...)
				assert.Nil(t, err)
				assert.Nil(t, resp)
				return nil
			},
		},
		{
			name: "GetCampaignsByIds need return none because not found tenant",
			auxFunc: func() error {
				campaignEntity := insertedValues[0]
				resp, err := repository.GetCampaignsByIds(context.Background(), campaignEntity.Id)
				assert.Nil(t, err)
				assert.Nil(t, resp)
				return err
			},
		},
		{
			name: "GetCampaignsTemplateById need return success",
			auxFunc: func() error {
				ids := make([]string, 0, len(insertedValues))
				responses := make(map[string]string)
				for _, campaignEntity := range insertedValues {
					ids = append(ids, campaignEntity.Id)
					responses[campaignEntity.Template.Id] = campaignEntity.Id
				}

				resp, err := repository.GetCampaignsTemplateById(ctx, ids)
				assert.Nil(t, err)
				assert.Equal(t, len(ids), len(resp))
				for _, id := range resp {
					if _, ok := responses[id]; !ok {
						return errors.New("not found")
					}
				}
				return err
			},
		},
		{
			name: "GetCampaignsTemplateById need return none because not found",
			auxFunc: func() error {
				resp, err := repository.GetCampaignsTemplateById(ctx, []string{"1"})
				assert.Nil(t, err)
				assert.Nil(t, resp)

				return err
			},
		},
		{
			name: "SaveCampaign need be success",
			auxFunc: func() error {
				name := "testSaveCampaign"
				campaignEntity := entity.Campaign{
					Name: name,
				}
				err := repository.SaveCampaign(ctx, &campaignEntity)
				assert.Nil(t, err)
				assert.Equal(t, name, campaignEntity.Name)
				assert.NotNil(t, campaignEntity.Id)

				return err
			},
		},
		{
			name:    "SaveCampaign need return error",
			isError: true,
			auxFunc: func() error {
				err := repository.SaveCampaign(ctx, nil)
				assert.NotNil(t, err)
				return err
			},
		},
		{
			name: "ListCampaing need return success",
			auxFunc: func() error {
				v, err := repository.ListCampaign(ctx)
				assert.Nil(t, err)
				assert.Equal(t, len(v), len(insertedValues))
				return err
			},
		},
		{
			name: "ListCampaing need return none because not found tenant",
			auxFunc: func() error {
				v, err := repository.ListCampaign(context.Background())
				assert.Nil(t, err)
				assert.Equal(t, len(v), 0)
				return err
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.auxFunc()
			if tt.isError {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func setupDb(ctx context.Context) []entity.Campaign {
	connection := hmongo.GetConnection(ctx)
	tenantId := hctx.Tenant.Get(ctx)
	v := []entity.Campaign{
		{
			Id:   primitive.NewObjectID().Hex(),
			Name: "test",
			Template: base.TemplateInfo{
				Id:   primitive.NewObjectID().Hex(),
				Name: "templateTest",
			},
		},
		{
			Id:   primitive.NewObjectID().Hex(),
			Name: "test1",
			Template: base.TemplateInfo{
				Id:   primitive.NewObjectID().Hex(),
				Name: "templateTest1",
			},
		},
	}

	docs := make([]interface{}, len(v))
	for i, d := range v {
		d.TenantId = tenantId
		docs[i] = d
	}

	_, err := connection.GetCollection(ctx, "campaigns").InsertMany(ctx, docs)
	if err != nil {
		panic(err)
	}
	return v
}
