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
	FindAll(ctx context.Context, logged *hctx.Logged) ([]domain.Template, error)
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

func (t *templateRepository) FindAll(ctx context.Context, logged *hctx.Logged) ([]domain.Template, error) {
	hlog.Debug(ctx, "templateRepository.FindAll", "finding all Templates")
	tenantId := hctx.Tenant.Get(ctx)
	cur, err := t.connection.GetCollection(ctx, t.collection).Find(ctx, bson.M{
		"tenantId": tenantId,
	})
	if err != nil {
		hlog.Error(ctx, "templateRepository.FindAll", err.Error())
		return nil, err
	}
	defer cur.Close(ctx)
	var templates []domain.Template
	for cur.Next(ctx) {
		var tmpl domain.Template
		if err = cur.Decode(&tmpl); err != nil {
			hlog.Error(ctx, "templateRepository.FindAll", err.Error())
			return nil, err
		}
		templates = append(templates, tmpl)
	}
	return templates, nil
}
