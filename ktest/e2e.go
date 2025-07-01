package ktest

import (
	"context"

	"github.com/inter-hubly/pilot/database/elasticsearch"
	"github.com/inter-hubly/pilot/database/hmongo"
	"github.com/inter-hubly/pilot/database/pgsql"
	"github.com/inter-hubly/pilot/testutils"
)

func MongoSetup(ctx context.Context) func(context.Context) error {
	mongoHost, close, err := testutils.Mongo(ctx)
	if err != nil {
		panic(err)
	}

	hmongo.NewConnection(
		ctx,
		hmongo.WithDatabase("test"),
		hmongo.WithUrl(mongoHost),
	)
	return close
}

func PgsqlSetup(ctx context.Context) func(context.Context) error {
	pgsqlHost, close, err := testutils.Pgsql(ctx)
	if err != nil {
		panic(err)
	}
	pgsql.NewConnection(
		ctx,
		pgsql.WithUrl(pgsqlHost),
	)
	return close
}

func ElasticSetup(ctx context.Context) func(context.Context) error {
	elasticHost, close, err := testutils.ElasticSearch(ctx)
	if err != nil {
		panic(err)
	}
	elasticsearch.NewConn(
		elasticsearch.WithUrl([]string{elasticHost}),
	)
	return close
}
