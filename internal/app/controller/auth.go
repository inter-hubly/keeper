package controller

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/internal/app/domain/dto"
	"github.com/inter-hubly/keeper/internal/app/service"
	"github.com/inter-hubly/keeper/internal/infraestructure/httprest"
)

type Auth interface {
	Login(c *gin.Context)
}

type authController struct {
	authService service.Authenticate
}

func NewAuth() *authController {
	var (
		authControllerOnce sync.Once
		auth               *authController
	)
	authControllerOnce.Do(func() {
		auth = &authController{
			authService: service.NewAuthenticate(),
		}
	})
	return auth
}

func (a *authController) Login(c *gin.Context) {
	var login dto.Login

	if err := c.BindJSON(&login); err != nil {
		httprest.Error(c, "Error when marshal login")
		return
	}

	accessToken, err := a.authService.Login(c, &login)
	if err != nil {
		httprest.Error(c, err.Error())
	}

	accessTokenDto := struct {
		AccessToken string `json:"accessToken"`
	}{
		AccessToken: accessToken,
	}

	httprest.Ok(c, accessTokenDto)
}
