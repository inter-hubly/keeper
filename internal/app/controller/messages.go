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

type Messages interface {
	SearchMessages(c *gin.Context)
}

type messagesController struct {
	messageService service.Messages
}

var (
	_messagesControllerOnce sync.Once
	_messagesController     *messagesController
)

func NewMessages(ctx context.Context) *messagesController {
	_messagesControllerOnce.Do(func() {
		_messagesController = &messagesController{
			messageService: service.NewMessages(ctx),
		}
	})
	return _messagesController
}

func (m *messagesController) SearchMessages(c *gin.Context) {
	var searchDTO *kdto.Search
	ctx, loggedUser := middleware.GetLoggedUser(c)

	allMsg, err := m.messageService.SearchMessages(ctx, loggedUser, searchDTO)
	if err != nil {
		httprest.Error(c, "Error find messages")
		return
	}
	httprest.Ok(c, allMsg)
}
