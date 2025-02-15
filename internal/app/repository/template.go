package repository

import (
	"context"
	"sync"

	"github.com/inter-hubly/keeper/internal/app/domain"
	"github.com/inter-hubly/pilot/database/hmongo"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Template interface {
	SaveTemplate(ctx context.Context, user *hctx.Logged, dto *domain.Template) (*domain.Template, error)
}

var (
	templateRepositoryOnce sync.Once
	template               *templateRepository
)

type templateRepository struct {
	connection hmongo.NoSqlConn
	collection string
}

func NewTemplate(ctx context.Context) *templateRepository {
	templateRepositoryOnce.Do(func() {
		template = &templateRepository{
			connection: hmongo.GetConnection(ctx),
			collection: "templates",
		}
	})
	return template
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
