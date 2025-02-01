package controller

import (
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

func NewMessages() *messagesController {
	var (
		messagesControllerOnce sync.Once
		messages               *messagesController
	)

	messagesControllerOnce.Do(func() {
		messages = &messagesController{
			messageService: service.NewMessages(),
		}
	})
	return messages
}

func (m *messagesController) SearchMessages(c *gin.Context) {
	var searchDTO *kdto.Search
	ctx, _ := middleware.GetLoggedUser(c)

	// if err := c.BindJSON(searchDTO); err != nil {
	// 	httprest.Error(c, "Error getting body")
	// 	return
	// }
	messages, err := m.messageService.SearchMessages(ctx, searchDTO)
	if err != nil {
		httprest.Error(c, "Error find messages")
		return
	}
	httprest.Ok(c, messages)
}
