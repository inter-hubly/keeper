package service

import (
	"context"
	"fmt"
	"github.com/inter-hubly/keeper/internal/app/domain"
	"sync"

	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
)

type Template interface {
	Save(ctx context.Context, user *hctx.Logged, dto domain.Template)
}

var (
	templateServiceOnce sync.Once
	template            *templateService
)

type templateService struct {
}

func NewTemplate(ctx context.Context) *templateService {
	templateServiceOnce.Do(func() {
		template = &templateService{}
	})
	return template
}

func (s *templateService) Save(ctx context.Context, user *hctx.Logged, dtoTemplate domain.Template) {
	hlog.Debug(ctx, "templateService.Save", fmt.Sprint("create template", dtoTemplate))
}
