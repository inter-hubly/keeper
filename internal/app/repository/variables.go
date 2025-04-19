//go:generate mockgen -source=variables.go -destination=mocks/variables_mock.go -package=mocks

package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/inter-hubly/keeper/internal/app/domain"
	"github.com/inter-hubly/pilot/database/hmongo"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Variable interface {
	SaveVariable(ctx context.Context, variable *domain.Variables) (*domain.Variables, error)
	SaveManyVariables(ctx context.Context, variable *domain.Variables) error
	GetVariables(ctx context.Context) (map[string]domain.SingleVariable, error)
}

type variableRepository struct {
	connection hmongo.NoSqlConn
	collection string
}

var (
	_variableRepositoryOnce sync.Once
	_variableRepository     *variableRepository
)

func NewVariables(ctx context.Context) *variableRepository {
	_variableRepositoryOnce.Do(func() {
		_variableRepository = &variableRepository{
			connection: hmongo.GetConnection(ctx),
			collection: "variables",
		}
	})
	return _variableRepository
}

func (r *variableRepository) SaveVariable(ctx context.Context, domainVariable *domain.Variables) (*domain.Variables, error) {
	hlog.Debug(ctx, "variableRepository.SaveVariable", "saving one variable")
	insertId, err := r.connection.GetCollection(ctx, r.collection).InsertOne(ctx, domainVariable)
	if err != nil {
		return nil, err
	}
	id := insertId.InsertedID
	domainVariable.Id = id.(primitive.ObjectID).Hex()
	return domainVariable, nil
}

func (r *variableRepository) SaveManyVariables(ctx context.Context, variables *domain.Variables) error {
	hlog.Debug(ctx, "variableRepository.SaveManyVariables", "saving many variable")
	tenantId := hctx.Tenant.Get(ctx)

	hasTenant, err := r.connection.GetCollection(ctx, r.collection).CountDocuments(ctx, bson.M{
		"tenantId": tenantId,
	})
	if err != nil {
		hlog.Error(ctx, "variableRepository.SaveManyVariables", "error counting tenant")
		return err
	}

	if hasTenant > 0 {
		if one, err := r.connection.GetCollection(ctx, r.collection).
			UpdateOne(ctx,
				bson.M{"tenantId": tenantId},
				bson.M{"$set": variables},
			); err == nil {
			if one.ModifiedCount == 0 {
				hlog.Error(ctx, "variableRepository.SaveManyVariables", "not modified variables")
				return errors.New("error updating variables")
			}
			return nil
		} else {
			hlog.Error(ctx, "variableRepository.SaveManyVariables", "error updating tenant")
			return err
		}
	}

	if _, err = r.connection.GetCollection(ctx, r.collection).InsertOne(ctx, variables); err != nil {
		hlog.Error(ctx, "variableRepository.SaveManyVariables", "error inserting variables")
		return err
	}

	return nil
}

func (r *variableRepository) GetVariables(ctx context.Context) (map[string]domain.SingleVariable, error) {
	hlog.Debug(ctx, "variableRepository.GetVariables", "getting variables for tenant")

	var result domain.Variables
	err := r.connection.GetCollection(ctx, r.collection).FindOne(ctx, bson.M{
		"tenantId": hctx.Tenant.Get(ctx),
	}).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// Caso o documento não seja encontrado, retornamos uma lista vazia, sem erro
			return map[string]domain.SingleVariable{}, nil
		}
		hlog.Error(ctx, "variableRepository.GetVariables", "error finding variables")
		return nil, err
	}

	if len(result.Variable) == 0 {
		return map[string]domain.SingleVariable{}, nil
	}

	return result.Variable, nil
}
