package messages

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/infraestructure/middleware"
	"github.com/inter-hubly/keeper/internal/domain/kdto"
	"github.com/inter-hubly/keeper/web/httprest"
)

type Controller interface {
	SearchMessages(c *gin.Context)
}

type messagesController struct {
	messageService Service
}

var (
	_messagesControllerOnce sync.Once
	_messagesController     *messagesController
)

func NewController(ctx context.Context) *messagesController {
	_messagesControllerOnce.Do(func() {
		_messagesController = &messagesController{
			messageService: NewService(ctx),
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
