package service

import (
	"context"
	"sync"
)

type Template interface {
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
