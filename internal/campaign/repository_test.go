package campaign

import (
	"context"
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

	repository := NewRepository(ctx)
	insertedValues := setupDb(ctx)

	for _, tt := range []struct {
		name    string
		isError bool
		auxFunc func() error
	}{
		{
			name: "GetCampaignById need return success",
			auxFunc: func() error {
				campaign := insertedValues[0]
				resp, err := repository.GetCampaignById(ctx, campaign.Id)
				assert.Nil(t, err)
				assert.Equal(t, campaign.Id, resp.Id)
				assert.Equal(t, campaign.Name, resp.Name)
				return nil
			},
		},
		{
			name: "GetCampaignById need return error",
		},
		{
			name: "GetCampaignById need return error without tenant",
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
		},
		{
			Id:   primitive.NewObjectID().Hex(),
			Name: "test1",
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
