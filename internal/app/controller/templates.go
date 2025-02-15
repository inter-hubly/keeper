package controller

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/internal/app/domain/kdto"
	"github.com/inter-hubly/keeper/internal/app/mediator"
	"github.com/inter-hubly/keeper/internal/infraestructure/httprest"
	"github.com/inter-hubly/keeper/internal/infraestructure/middleware"
)

type Templates interface {
	Save(c *gin.Context)
}

var (
	templatesControllerOnce sync.Once
	templates               *templateController
)

type templateController struct {
	templateMediator mediator.Template
}

func NewTemplate(ctx context.Context) *templateController {
	templatesControllerOnce.Do(func() {
		templates = &templateController{
			templateMediator: mediator.NewTemplate(ctx),
		}
	})
	return templates
}

func (t *templateController) Save(c *gin.Context) {
	var templateDto kdto.Template
	ctx, loggedUser := middleware.GetLoggedUser(c)
	if err := c.BindJSON(&templateDto); err != nil {
		httprest.Error(c, "Error when marshal body")
		return
	}

	t.templateMediator.Save(ctx, loggedUser, templateDto)
}
