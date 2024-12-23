package express

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/inter-hubly/pilot/database/elasticsearch"
	"github.com/inter-hubly/pilot/database/pgsql"
	"github.com/inter-hubly/pilot/server"
)

func Start(engine *gin.Engine) {
	engine.Use(corsMiddleware())
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

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Configurar os cabeçalhos CORS
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
