package service

import (
	"context"
	"sync"

	"github.com/inter-hubly/keeper/internal/app/domain"
	"github.com/inter-hubly/keeper/internal/app/domain/kdto"
	"github.com/inter-hubly/keeper/internal/app/repository"
	"github.com/inter-hubly/pilot/domain/base"
	"github.com/inter-hubly/pilot/domain/valueobject"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
)

type Variables interface {
	SaveVariables(ctx context.Context, logged *hctx.Logged, variableDto *kdto.Variable) error
	SaveManyVariables(ctx context.Context, logged *hctx.Logged, variableDto []kdto.Variable) error
	SearchVariables(ctx context.Context) ([]kdto.Variable, error)
}

type variablesService struct {
	variableRepository repository.Variable
}

var (
	_variablesServiceOnce sync.Once
	_variablesService     *variablesService
)

func NewVariables(ctx context.Context) *variablesService {
	_variablesServiceOnce.Do(func() {
		_variablesService = &variablesService{
			variableRepository: repository.NewVariables(ctx),
		}
	})
	return _variablesService
}

func (v *variablesService) SaveVariables(ctx context.Context, logged *hctx.Logged, variableDto *kdto.Variable) error {
	slug := valueobject.NewSlug(variableDto.Label)
	variableMap := make(map[string]domain.SingleVariable)
	variableMap[slug.Value()] = domain.SingleVariable{
		Slug:  slug.Value(),
		Label: variableDto.Label,
		Type:  variableDto.Type,
	}

	newVariable := domain.Variables{
		Variable: variableMap,
		Entity:   base.NewBaseEntity(ctx, logged),
	}

	if _, err := v.variableRepository.SaveVariable(ctx, &newVariable); err != nil {
		return err
	}
	return nil
}

func (v *variablesService) SaveManyVariables(ctx context.Context, logged *hctx.Logged, variableDto []kdto.Variable) error {
	hlog.Debug(ctx, "variablesService.SaveManyVariables", "saving many variables")
	manyVariables := map[string]domain.SingleVariable{}
	for _, vdto := range variableDto {
		slug := valueobject.NewSlug(vdto.Label)

		newVariable := domain.SingleVariable{
			Slug:  slug.Value(),
			Label: vdto.Label,
			Type:  vdto.Type,
		}
		manyVariables[slug.Value()] = newVariable
	}
	variablesDb := domain.Variables{
		Variable: manyVariables,
		Entity:   base.NewBaseEntity(ctx, logged),
	}
	if err := v.variableRepository.SaveManyVariables(ctx, &variablesDb); err != nil {
		hlog.Error(ctx, "variablesService.SaveManyVariables", err.Error())
		return err
	}
	return nil
}

func (v *variablesService) SearchVariables(ctx context.Context) ([]kdto.Variable, error) {
	hlog.Debug(ctx, "variablesService.GetVariables", "getting variables")
	getVariables, err := v.variableRepository.GetVariables(ctx)
	if err != nil {
		hlog.Error(ctx, "variablesService.GetVariables", err.Error())
		return nil, err
	}
	resp := make([]kdto.Variable, 0, len(getVariables))
	for _, vdto := range getVariables {
		resp = append(resp, kdto.Variable{
			Label: vdto.Label,
			Type:  vdto.Type,
			Slug:  vdto.Slug,
		})
	}
	return resp, nil
}
