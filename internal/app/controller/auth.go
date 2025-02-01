package controller

import (
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

func NewAuth() *authController {

	authControllerOnce.Do(func() {
		auth = &authController{
			authService: service.NewAuthenticate(),
		}
	})
	return auth
}

func (a *authController) Login(c *gin.Context) {
	var login kdto.Login

	if err := c.BindJSON(&login); err != nil {
		httprest.Error(c, "Error when marshal login")
		return
	}

	accessToken, err := a.authService.Login(c, &login)
	if err != nil {
		httprest.Unauthorized(c)
		return
	}

	accessTokenDto := struct {
		AccessToken string `json:"accessToken"`
	}{
		AccessToken: accessToken,
	}

	httprest.Ok(c, accessTokenDto)
}

func (a *authController) CreateUser(c *gin.Context) {
	var user kdto.User

	if err := c.BindJSON(&user); err != nil {
		httprest.Error(c, "Error when marshal login")
		return
	}

	if err := a.authService.CreateUser(c, &user); err != nil {
		httprest.Error(c, err.Error())
		return
	}

	httprest.Ok(c, nil)
}
