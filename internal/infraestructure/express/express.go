package express

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	rabbitmq "github.com/inter-hubly/pilot/broker"
	"github.com/inter-hubly/pilot/database/elasticsearch"
	"github.com/inter-hubly/pilot/database/hmongo"
	"github.com/inter-hubly/pilot/database/pgsql"
	"github.com/inter-hubly/pilot/server"
)

const ExchangeBroker = "linker"

func Start(ctx context.Context, engine *gin.Engine) {
	engine.Use(corsMiddleware())
	pgsql.NewConnection(
		pgsql.WithUrl(server.GetPgsqlConfig().Host),
	)

	hmongo.NewConnection(
		ctx,
		hmongo.WithDatabase(server.GetMongoConfig().Database),
		hmongo.WithUrl(server.GetMongoConfig().Host),
	)

	elasticsearch.NewConn(
		elasticsearch.WithUrl([]string{server.GetElasticSearch().Host}),
		elasticsearch.WithUsernameAndPassword(
			server.GetElasticSearch().Username,
			server.GetElasticSearch().Password,
		),
	)
	rabbitmq.NewRabbitMQ(ctx, ExchangeBroker, "topic", rabbitmq.WithURL(server.GetAmpqConfig().Host))
	if err := rabbitmq.GetConnection().
		QueueBind(
			ctx,
			rabbitmq.NewQueueBinding("campaign.init", "campaign.init", ExchangeBroker),
		); err != nil {
		panic(err)
	}

	NewKeeperController(ctx, engine)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, tenant")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
