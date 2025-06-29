//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go -package=mocks

package templates

import (
	"context"
	"sync"

	"github.com/inter-hubly/keeper/internal/domain"
	"github.com/inter-hubly/pilot/database/hmongo"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	SaveTemplate(ctx context.Context, dto *domain.Template) (*domain.Template, error)
	SearchTemplates(ctx context.Context, logged *hctx.Logged) ([]domain.Template, error)
	GetTemplateByIds(ctx context.Context, ids []string) ([]domain.Template, error)
	GetTemplateById(ctx context.Context, id string) (*domain.Template, error)
	SaveManyTemplate(ctx context.Context, save []domain.Template) error
}

var (
	_templateRepositoryOnce sync.Once
	_templateRepository     *templateRepository
)

type templateRepository struct {
	connection hmongo.NoSqlConn
	collection string
}

func NewRepository(ctx context.Context) *templateRepository {
	_templateRepositoryOnce.Do(func() {
		_templateRepository = &templateRepository{
			connection: hmongo.GetConnection(ctx),
			collection: "templates",
		}
	})
	return _templateRepository
}

func (t *templateRepository) SaveTemplate(ctx context.Context, domainTemplate *domain.Template) (*domain.Template, error) {
	hlog.Debug(ctx, "template.repository.SaveTemplate", "saving one Template")

	primitiveObjectId, err := primitive.ObjectIDFromHex(domainTemplate.Id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"tenantId": domainTemplate.TenantId,
		"_id":      primitiveObjectId,
	}

	// preciso retirar o _id
	domainTemplate.Id = ""
	update := bson.M{
		"$set": domainTemplate,
	}

	opts := options.Update().SetUpsert(true)

	result, err := t.connection.GetCollection(ctx, t.collection).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}
	if result.UpsertedID != nil {
		domainTemplate.Id = result.UpsertedID.(primitive.ObjectID).Hex()
	}

	return domainTemplate, nil
}

func (t *templateRepository) SearchTemplates(ctx context.Context, logged *hctx.Logged) ([]domain.Template, error) {
	hlog.Debug(ctx, "template.repository.SearchTemplates", "finding all Templates")
	tenantId := hctx.Tenant.Get(ctx)
	cur, err := t.connection.GetCollection(ctx, t.collection).Find(ctx, bson.M{
		"tenantId": tenantId,
	})
	if err != nil {
		hlog.Error(ctx, "template.repository.SearchTemplates", err.Error())
		return nil, err
	}
	defer cur.Close(ctx)
	var templates []domain.Template
	if err = cur.All(ctx, &templates); err != nil {
		hlog.Error(ctx, "template.repository.SearchTemplates", err.Error())
		return nil, err
	}

	return templates, nil
}

func (t *templateRepository) GetTemplateByIds(ctx context.Context, ids []string) ([]domain.Template, error) {
	hlog.Debug(ctx, "template.repository.GetTemplateByIds", "finding all Templates")
	tenantId := hctx.Tenant.Get(ctx)
	var objectIds []primitive.ObjectID
	for _, id := range ids {
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			hlog.Error(ctx, "template.repository.GetTemplateByIds", err.Error())
			return nil, err
		}

		objectIds = append(objectIds, objectId)
	}
	cur, err := t.connection.GetCollection(ctx, t.collection).Find(ctx, bson.M{
		"tenantId": tenantId,
		"_id": bson.M{
			"$in": objectIds,
		},
	})
	if err != nil {
		hlog.Error(ctx, "template.repository.GetTemplateByIds", err.Error())
		return nil, err
	}
	defer cur.Close(ctx)
	var templates []domain.Template
	if err = cur.All(ctx, &templates); err != nil {
		hlog.Error(ctx, "template.repository.GetTemplateByIds", err.Error())
		return nil, err
	}
	return templates, nil
}

func (t *templateRepository) GetTemplateById(ctx context.Context, id string) (*domain.Template, error) {
	hlog.Debug(ctx, "template.repository.GetTemplateById", "finding all Templates")
	tenantId := hctx.Tenant.Get(ctx)
	var (
		objectId primitive.ObjectID
		err      error
	)

	objectId, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		hlog.Error(ctx, "template.repository.GetTemplateByIds", err.Error())
		return nil, err
	}

	var template domain.Template
	if err = t.connection.
		GetCollection(ctx, t.collection).
		FindOne(ctx, bson.M{
			"_id":      objectId,
			"tenantId": tenantId,
		}).
		Decode(&template); err != nil {
		hlog.Error(ctx, "template.repository.GetTemplateById", "error find template")
	}
	return &template, nil
}

func (t *templateRepository) SaveManyTemplate(ctx context.Context, save []domain.Template) error {
	hlog.Debug(ctx, "template.repository.SaveManyTemplate", "saving all Templates")
	tenantId := hctx.Tenant.Get(ctx)

	var operations []mongo.WriteModel

	docs := make([]interface{}, len(save))
	for i := range save {
		filter := bson.M{
			"slug":     save[i].Slug,
			"tenantId": tenantId,
		}

		save[i].TenantId = tenantId
		docs[i] = save[i]

		update := bson.M{
			"$set": save[i],
		}
		model := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)
		operations = append(operations, model)
	}
	if len(operations) == 0 {
		return nil
	}
	_, err := t.connection.GetCollection(ctx, t.collection).BulkWrite(ctx, operations)
	if err != nil {
		hlog.Error(ctx, "template.repository.SaveManyTemplate", err.Error())
		return err
	}
	return nil
}
