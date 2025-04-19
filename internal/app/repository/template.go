//go:generate mockgen -source=template.go -destination=mocks/template_mock.go -package=mocks

package repository

import (
	"context"
	"sync"

	"github.com/inter-hubly/keeper/internal/app/domain"
	"github.com/inter-hubly/pilot/database/hmongo"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Template interface {
	SaveTemplate(ctx context.Context, user *hctx.Logged, dto *domain.Template) (*domain.Template, error)
	SearchTemplates(ctx context.Context, logged *hctx.Logged) ([]domain.Template, error)
	GetTemplateByIds(ctx context.Context, ids []string) ([]domain.Template, error)
	GetTemplateById(ctx context.Context, id string) (*domain.Template, error)
}

var (
	_templateRepositoryOnce sync.Once
	_template               *templateRepository
)

type templateRepository struct {
	connection hmongo.NoSqlConn
	collection string
}

func NewTemplate(ctx context.Context) *templateRepository {
	_templateRepositoryOnce.Do(func() {
		_template = &templateRepository{
			connection: hmongo.GetConnection(ctx),
			collection: "templates",
		}
	})
	return _template
}

func (t *templateRepository) SaveTemplate(ctx context.Context, user *hctx.Logged, domainTemplate *domain.Template) (*domain.Template, error) {
	hlog.Debug(ctx, "templateRepository.SaveTemplate", "saving one Template")
	insertId, err := t.connection.GetCollection(ctx, t.collection).InsertOne(ctx, domainTemplate)
	if err != nil {
		return nil, err
	}
	id := insertId.InsertedID
	domainTemplate.Id = id.(primitive.ObjectID).Hex()
	return domainTemplate, nil
}

func (t *templateRepository) SearchTemplates(ctx context.Context, logged *hctx.Logged) ([]domain.Template, error) {
	hlog.Debug(ctx, "templateRepository.SearchTemplates", "finding all Templates")
	tenantId := hctx.Tenant.Get(ctx)
	cur, err := t.connection.GetCollection(ctx, t.collection).Find(ctx, bson.M{
		"tenantId": tenantId,
	})
	if err != nil {
		hlog.Error(ctx, "templateRepository.SearchTemplates", err.Error())
		return nil, err
	}
	defer cur.Close(ctx)
	var templates []domain.Template
	if err = cur.All(ctx, &templates); err != nil {
		hlog.Error(ctx, "templateRepository.SearchTemplates", err.Error())
		return nil, err
	}

	return templates, nil
}

func (t *templateRepository) GetTemplateByIds(ctx context.Context, ids []string) ([]domain.Template, error) {
	hlog.Debug(ctx, "templateRepository.GetTemplateByIds", "finding all Templates")
	tenantId := hctx.Tenant.Get(ctx)
	var objectIds []primitive.ObjectID
	for _, id := range ids {
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			hlog.Error(ctx, "templateRepository.GetTemplateByIds", err.Error())
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
		hlog.Error(ctx, "templateRepository.GetTemplateByIds", err.Error())
		return nil, err
	}
	defer cur.Close(ctx)
	var templates []domain.Template
	if err = cur.All(ctx, &templates); err != nil {
		hlog.Error(ctx, "templateRepository.GetTemplateByIds", err.Error())
		return nil, err
	}
	return templates, nil
}

func (t *templateRepository) GetTemplateById(ctx context.Context, id string) (*domain.Template, error) {
	hlog.Debug(ctx, "templateRepository.GetTemplateById", "finding all Templates")
	tenantId := hctx.Tenant.Get(ctx)
	var (
		objectId primitive.ObjectID
		err      error
	)

	objectId, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		hlog.Error(ctx, "templateRepository.GetTemplateByIds", err.Error())
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
		hlog.Error(ctx, "templateRepository.GetTemplateById", "error find template")
	}
	return &template, nil
}
