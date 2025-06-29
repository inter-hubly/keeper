package variables

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/infraestructure/middleware"
	"github.com/inter-hubly/keeper/internal/domain/kdto"
	"github.com/inter-hubly/keeper/web/httprest"
)

type Controller interface {
	SearchVariables(c *gin.Context)
	SaveVariable(c *gin.Context)
	SaveManyVariable(c *gin.Context)
}

var (
	_variableControllerOnce sync.Once
	_variableController     *variableController
)

type variableController struct {
	variablesService Service
}

func NewController(ctx context.Context) *variableController {
	_variableControllerOnce.Do(func() {
		_variableController = &variableController{
			variablesService: NewService(ctx),
		}
	})
	return _variableController
}

func (v *variableController) SaveVariable(c *gin.Context) {
	ctx, loggedUser := middleware.GetLoggedUser(c)

	var variableDto kdto.Variable

	if err := c.BindJSON(&variableDto); err != nil {
		httprest.Error(c, "Error when transform body")
		return
	}

	if err := v.variablesService.SaveVariables(ctx, loggedUser, &variableDto); err != nil {
		httprest.Error(c, "Error when save variable")
		return
	}
	httprest.Created(c, nil)
}

func (v *variableController) SaveManyVariable(c *gin.Context) {
	ctx, loggedUser := middleware.GetLoggedUser(c)

	var variableDto []kdto.Variable

	if err := c.BindJSON(&variableDto); err != nil {
		httprest.Error(c, "Error when marshal body")
		return
	}

	if err := v.variablesService.SaveManyVariables(ctx, loggedUser, variableDto); err != nil {
		httprest.Error(c, "Error when save variable")
		return
	}
	httprest.Created(c, nil)
}

func (v *variableController) SearchVariables(c *gin.Context) {
	ctx, _ := middleware.GetLoggedUser(c)

	getVariables, err := v.variablesService.SearchVariables(ctx)
	if err != nil {
		httprest.Error(c, "Error when get variables")
		return
	}
	httprest.Ok(c, getVariables)
}
