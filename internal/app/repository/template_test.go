package repository

import (
	"context"
	"testing"

	"github.com/inter-hubly/keeper/internal/app/domain"
	"github.com/inter-hubly/pilot/database/hmongo"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/testutils"
	"github.com/stretchr/testify/assert"
)

func TestFindManyContact(t *testing.T) {
	ctx := testutils.SetLoggedUser(context.Background())
	host, close, err := testutils.Mongo(ctx)
	if err != nil {
		panic(err)
	}
	defer close(ctx)

	if hmongo.GetConnection(ctx) == nil {
		hmongo.NewConnection(
			ctx,
			hmongo.WithDatabase("test"),
			hmongo.WithUrl(host),
		)
	}

	repository := templateRepository{
		connection: hmongo.GetConnection(ctx),
		collection: "template",
	}

	for _, v := range []struct {
		testName string
		auxFunc  func()
	}{
		{
			testName: "Need to save many templates",
			auxFunc: func() {
				allTemplates := []domain.Template{
					{
						Name: "test1",
						Slug: "test1",
					},
					{
						Name: "test2",
						Slug: "test2",
					},
				}

				err = repository.SaveManyTemplate(ctx, allTemplates)
				assert.NoError(t, err)
				logged := hctx.LoggedUser.Get(ctx)
				templates, err := repository.SearchTemplates(ctx, &logged)
				assert.NoError(t, err)
				for i := range templates {
					assert.Equal(t, allTemplates[i].Name, templates[i].Name)
					assert.Equal(t, allTemplates[i].TenantId, logged.Tenant)
				}
			},
		},
	} {
		t.Run(v.testName, func(t *testing.T) {
			v.auxFunc()
		})
	}
}
