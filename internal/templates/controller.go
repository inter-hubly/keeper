package templates

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/infraestructure/middleware"
	"github.com/inter-hubly/keeper/internal/domain"
	"github.com/inter-hubly/keeper/internal/domain/kdto"
	"github.com/inter-hubly/keeper/web/httprest"
)

type Controller interface {
	Save(c *gin.Context)
	SearchTemplates(c *gin.Context)
	SincronizeWhatsAppTemplate(c *gin.Context)
	SaveVariables(c *gin.Context)
	CountVariables(c *gin.Context)
}

var (
	_templateControllerOnce sync.Once
	_templateController     *templateController
)

type templateController struct {
	templateMediator Mediator
	templateService  Service
}

func NewController(ctx context.Context) *templateController {
	_templateControllerOnce.Do(func() {
		_templateController = &templateController{
			templateMediator: NewMediator(ctx),
			templateService:  NewService(ctx),
		}
	})
	return _templateController
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

	// TODO adicionar contexto de tempo
	go func() {
		t.templateMediator.Save(ctx, loggedUser, &templateDto)
	}()

	httprest.Created(c, nil)
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

func (t *templateController) SincronizeWhatsAppTemplate(c *gin.Context) {
	ctx, loggedUser := middleware.GetLoggedUser(c)
	if err := t.templateService.SincronizeWhatsAppTemplate(ctx, loggedUser); err != nil {
		httprest.Error(c, "Error when sincronize template")
		return
	}
	httprest.Created(c, nil)
}

func (t *templateController) SaveVariables(c *gin.Context) {
	templateId := c.Param("templateId")
	ctx, _ := middleware.GetLoggedUser(c)

	var variables []kdto.Variable

	if err := c.BindJSON(&variables); err != nil {
		httprest.Error(c, "Error when marshal body")
		return
	}

	if err := t.templateService.SaveVariable(ctx, variables, templateId); err != nil {
		httprest.Error(c, "Error when save variables")
		return
	}
	httprest.Created(c, nil)
}

func (t *templateController) CountVariables(c *gin.Context) {
	ctx, _ := middleware.GetLoggedUser(c)
	templateId := c.Param("templateId")
	if templateId == "" {
		httprest.Error(c, "Error when get template id from request")
		return
	}

	variables, err := t.templateService.CountVariables(ctx, templateId)
	if err != nil {
		httprest.Error(c, "Error when count variables")
	}
	httprest.Ok(c, variables)
}
