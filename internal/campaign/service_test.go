package campaign_test

import (
	"context"
	"github.com/inter-hubly/keeper/internal/campaign"
	"github.com/inter-hubly/keeper/internal/domain/kdto"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService(t *testing.T) {
	ctx := testutils.SetLoggedUser(context.Background())
	service := campaign.NewService(ctx)
	assert.NotNil(t, service)
	for _, tt := range []struct {
		name    string
		isError bool
		auxFunc func() error
	}{
		{
			name: "SaveCampaign need be return success",
			auxFunc: func() error {
				loggedUser := hctx.LoggedUser.Get(ctx)
				dto := getCampaignDto(ctx)
				saveCampaign, err := service.SaveCampaign(ctx, &loggedUser, &dto)
				assert.Nil(t, err)
				assert.NotNil(t, saveCampaign)
				return nil
			},
		},
		{
			name:    "SaveCampaign need be return error for if 1",
			isError: true,
			auxFunc: func() error {
				// TODO: implement test for SaveCampaign (if 1)
				return nil
			},
		},
		{
			name:    "SaveCampaign need be return error for if 2",
			isError: true,
			auxFunc: func() error {
				// TODO: implement test for SaveCampaign (if 2)
				return nil
			},
		},
		{
			name:    "SaveCampaign need be return error for if 3",
			isError: true,
			auxFunc: func() error {
				// TODO: implement test for SaveCampaign (if 3)
				return nil
			},
		},
		{
			name:    "SaveCampaign need be return error for if 4",
			isError: true,
			auxFunc: func() error {
				// TODO: implement test for SaveCampaign (if 4)
				return nil
			},
		},
		{
			name:    "SaveCampaign need be return error for if 5",
			isError: true,
			auxFunc: func() error {
				// TODO: implement test for SaveCampaign (if 5)
				return nil
			},
		},
		{
			name:    "SaveCampaign need be return error for if 6",
			isError: true,
			auxFunc: func() error {
				// TODO: implement test for SaveCampaign (if 6)
				return nil
			},
		},

		{
			name: "GetCampaign need be return success",
			auxFunc: func() error {
				// TODO: implement test for GetCampaign
				return nil
			},
		},
		{
			name:    "GetCampaign need be return error for if 1",
			isError: true,
			auxFunc: func() error {
				// TODO: implement test for GetCampaign (if 1)
				return nil
			},
		},

		{
			name: "StartCampaign need be return success",
			auxFunc: func() error {
				// TODO: implement test for StartCampaign
				return nil
			},
		},
		{
			name:    "StartCampaign need be return error for if 1",
			isError: true,
			auxFunc: func() error {
				// TODO: implement test for StartCampaign (if 1)
				return nil
			},
		},

		{
			name: "ListCampaign need be return success",
			auxFunc: func() error {
				// TODO: implement test for ListCampaign
				return nil
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

func getCampaignDto(ctx context.Context) kdto.Campaign {
	return kdto.Campaign{
		Name:       "Test Campaign",
		TemplateId: "testTemplateId",
	}
}
