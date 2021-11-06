package schema

import (
	"github.com/Ekenzy-101/Go-GraphQL-API/graphql/query"
	"github.com/graphql-go/graphql"
)

var Query = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"me":    query.Me,
		"post":  query.Post,
		"posts": query.Posts,
	},
})
