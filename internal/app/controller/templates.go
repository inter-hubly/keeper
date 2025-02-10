package controller

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/internal/app/service"
)

type Templates interface {
	Get(c *gin.Context)
}

var (
	templatesControllerOnce sync.Once
	templates               *templateController
)

type templateController struct {
	templateService service.Template
}

func NewTemplate(ctx context.Context) *templateController {
	templatesControllerOnce.Do(func() {
		templates = &templateController{
			templateService: service.NewTemplate(ctx),
		}
	})
	return templates
}

func (t *templateController) Get(c *gin.Context) {

}
