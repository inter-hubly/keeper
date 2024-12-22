package express

import (
	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/pilot/database/elasticsearch"
	"github.com/inter-hubly/pilot/database/pgsql"
	"github.com/inter-hubly/pilot/server"
)

func Start(engine *gin.Engine) {
	pgsql.NewConnection(
		pgsql.WithUrl(server.GetPgsqlConfig().Host),
	)

	elasticsearch.NewConn(
		elasticsearch.WithUrl([]string{server.GetElasticSearch().Host}),
		elasticsearch.WithUsernameAndPassword(
			server.GetElasticSearch().Username,
			server.GetElasticSearch().Password,
		),
	)

	NewKeeperController(engine)
}
