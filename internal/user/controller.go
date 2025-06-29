package user

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/internal/domain/kdto"
	"github.com/inter-hubly/keeper/web/httprest"
)

type Controller interface {
	Login(c *gin.Context)
	CreateUser(c *gin.Context)
}

type authController struct {
	authService Service
}

var (
	_userControllerOnce sync.Once
	_userController     *authController
)

func NewController(ctx context.Context) *authController {
	_userControllerOnce.Do(func() {
		_userController = &authController{
			authService: NewAuthenticate(ctx),
		}
	})
	return _userController
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
