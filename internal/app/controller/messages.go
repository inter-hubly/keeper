package controller

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/internal/app/domain/dto"
	"github.com/inter-hubly/keeper/internal/app/service"
	"github.com/inter-hubly/keeper/internal/infraestructure/httprest"
	"github.com/inter-hubly/pilot/hctx"
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
	var searchDTO *dto.Search
	tenantHeader := c.GetHeader("tenant")
	ctx := hctx.Tenant.New(tenantHeader)

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
