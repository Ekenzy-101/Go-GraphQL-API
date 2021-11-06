package schema

import (
	"github.com/graphql-go/graphql"
)

func New() (graphql.Schema, error) {
	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    Query,
		Mutation: Mutation,
	})
}
