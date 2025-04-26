package controller

import (
	"context"
	"sync"

	"github.com/inter-hubly/keeper/internal/app/domain"
	"github.com/inter-hubly/keeper/internal/app/service"
	"github.com/inter-hubly/pilot/hlog"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/internal/app/mediator"
	"github.com/inter-hubly/keeper/internal/infraestructure/httprest"
	"github.com/inter-hubly/keeper/internal/infraestructure/middleware"
)

type Templates interface {
	Save(c *gin.Context)
	SearchTemplates(c *gin.Context)
	SincronizeWhatsAppTemplate(c *gin.Context)
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
	hlog.Debug(c, "templateController.Save", "Save Template")
	var templateDto domain.Template
	ctx, loggedUser := middleware.GetLoggedUser(c)
	if err := c.BindJSON(&templateDto); err != nil {
		httprest.Error(c, "Error when marshal body")
		return
	}
	if templateDto.Language == "" {
		templateDto.Language = "pt_BR"
	}

	// TODO adicionar contexto de tempo
	go func() {
		t.templateMediator.Save(ctx, loggedUser, &templateDto)
	}()

	httprest.Created(c, nil)
}

func (t *templateController) SearchTemplates(c *gin.Context) {
	hlog.Debug(c, "templateController.SearchTemplates", "SearchTemplates Template")

	ctx, loggedUser := middleware.GetLoggedUser(c)
	all, err := t.templateService.SearchTemplates(ctx, loggedUser)
	if err != nil {
		httprest.Error(c, "Error when find all templates")
		return
	}
	httprest.Ok(c, all)
}

func (t *templateController) SincronizeWhatsAppTemplate(c *gin.Context) {
	hlog.Debug(c, "templateController.SincronizeWhatsAppTemplate", "Sincronize Template")

	ctx, loggedUser := middleware.GetLoggedUser(c)
	if err := t.templateService.SincronizeWhatsAppTemplate(ctx, loggedUser); err != nil {
		httprest.Error(c, "Error when sincronize template")
		return
	}
	httprest.Created(c, nil)
}
