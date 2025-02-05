package controller

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/internal/app/domain/kdto"
	"github.com/inter-hubly/keeper/internal/app/service"
	"github.com/inter-hubly/keeper/internal/infraestructure/httprest"
)

type Auth interface {
	Login(c *gin.Context)
	CreateUser(c *gin.Context)
}

type authController struct {
	authService service.Authenticate
}

var (
	authControllerOnce sync.Once
	auth               *authController
)

func NewAuth(ctx context.Context) *authController {

	authControllerOnce.Do(func() {
		auth = &authController{
			authService: service.NewAuthenticate(ctx),
		}
	})
	return auth
}

func (a *authController) Login(c *gin.Context) {
	var login kdto.Login

	if err := c.BindJSON(&login); err != nil {
		httprest.Error(c, "Error when marshal body")
		return
	}

	accessToken, err := a.authService.Login(c, &login)
	if err != nil {
		httprest.Unauthorized(c)
		return
	}

	httprest.Ok(c, accessToken)
}

func (a *authController) CreateUser(c *gin.Context) {
	var user kdto.User

	if err := c.BindJSON(&user); err != nil {
		httprest.Error(c, "Error when marshal body")
		return
	}

	if err := a.authService.CreateUser(c, &user); err != nil {
		httprest.Error(c, err.Error())
		return
	}

	httprest.Ok(c, nil)
}
