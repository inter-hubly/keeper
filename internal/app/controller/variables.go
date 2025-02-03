package controller

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/internal/app/domain/kdto"
	"github.com/inter-hubly/keeper/internal/app/service"
	"github.com/inter-hubly/keeper/internal/infraestructure/httprest"
	"github.com/inter-hubly/keeper/internal/infraestructure/middleware"
)

type Variables interface {
	GetVariables(c *gin.Context)
	SaveVariable(c *gin.Context)
	SaveManyVariable(c *gin.Context)
}

var (
	variableControllerOnce sync.Once
	variables              *variableController
)

type variableController struct {
	variablesService service.Variables
}

func NewVariable(ctx context.Context) *variableController {

	variableControllerOnce.Do(func() {
		variables = &variableController{
			variablesService: service.NewVariables(ctx),
		}
	})
	return variables
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

func (v *variableController) GetVariables(c *gin.Context) {
	ctx, _ := middleware.GetLoggedUser(c)

	getVariables, err := v.variablesService.GetVariables(ctx)
	if err != nil {
		httprest.Error(c, "Error when get variables")
		return
	}
	httprest.Ok(c, getVariables)
}
