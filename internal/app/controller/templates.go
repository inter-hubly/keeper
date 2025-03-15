package controller

import (
	"context"
	"sync"

	"github.com/inter-hubly/keeper/internal/app/domain"
	"github.com/inter-hubly/keeper/internal/app/service"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/internal/app/mediator"
	"github.com/inter-hubly/keeper/internal/infraestructure/httprest"
	"github.com/inter-hubly/keeper/internal/infraestructure/middleware"
)

type Templates interface {
	Save(c *gin.Context)
	SearchTemplates(c *gin.Context)
}

var (
	templatesControllerOnce sync.Once
	templates               *templateController
)

type templateController struct {
	templateMediator mediator.Template
	templateService  service.Template
}

func NewTemplate(ctx context.Context) *templateController {
	templatesControllerOnce.Do(func() {
		templates = &templateController{
			templateMediator: mediator.NewTemplate(ctx),
			templateService:  service.NewTemplate(ctx),
		}
	})
	return templates
}

func (t *templateController) Save(c *gin.Context) {
	var templateDto domain.Template
	ctx, loggedUser := middleware.GetLoggedUser(c)
	if err := c.BindJSON(&templateDto); err != nil {
		httprest.Error(c, "Error when marshal body")
		return
	}
	if templateDto.Language == "" {
		templateDto.Language = "pt_BR"
	}

	savedValue, err := t.templateMediator.Save(ctx, loggedUser, &templateDto)
	if err != nil {
		httprest.Error(c, "Error when save template")
		return
	}
	httprest.Created(c, savedValue)
}

func (t *templateController) SearchTemplates(c *gin.Context) {
	ctx, loggedUser := middleware.GetLoggedUser(c)
	all, err := t.templateService.SearchTemplates(ctx, loggedUser)
	if err != nil {
		httprest.Error(c, "Error when find all templates")
		return
	}
	httprest.Ok(c, all)
}
