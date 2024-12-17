package express

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/keeper/internal/app/controller"
)

var (
	controllersOnce   sync.Once
	keeperControllers *controllers
)

type controllers struct {
	engine           *gin.Engine
	clientController controller.Client
}

func NewKeeperController(engine *gin.Engine) {
	controllersOnce.Do(func() {
		keeperControllers = &controllers{
			engine:           engine,
			clientController: controller.NewClient(),
		}
	})
	keeperControllers.startControllers()
}

func (c *controllers) startControllers() {
	{
		clientGroup := c.engine.Group("/client")
		clientGroup.GET("/:id", c.clientController.GetClient)
	}

}
