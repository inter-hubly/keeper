package infraestructure

import (
	"context"

	"github.com/inter-hubly/keeper/web/router"
	rabbitmq "github.com/inter-hubly/pilot/broker"
	"github.com/inter-hubly/pilot/database/elasticsearch"
	"github.com/inter-hubly/pilot/database/hmongo"
	"github.com/inter-hubly/pilot/database/pgsql"
	"github.com/inter-hubly/pilot/server"
)

const ExchangeBroker = "linker"

func Start(ctx context.Context) {
	pgsql.NewConnection(
		ctx,
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

	router.NewRouter(ctx)
}
