package query

import (
	"errors"

	"github.com/Ekenzy-101/Go-GraphQL-API/config"
	"github.com/Ekenzy-101/Go-GraphQL-API/graphql/types"
	"github.com/Ekenzy-101/Go-GraphQL-API/service"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

var Posts = &graphql.Field{
	Type:        graphql.NewList(types.Post),
	Description: "Get the latest posts",
	Args:        types.PaginationArgs,
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		return func() (interface{}, error) {
			skip, ok := params.Args["skip"].(uint64)
			if !ok {
				skip = 0
			}

			limit, ok := params.Args["limit"].(uint64)
			if !ok {
				limit = 2
			}

			appService := params.Context.Value(config.ServiceKey).(service.Service)
			posts, err := appService.GetLatestPosts(params.Context, map[string]uint64{"skip": skip, "limit": limit})
			if err != nil {
				return nil, err
			}
			return posts, nil
		}, nil
	},
}

var Post = &graphql.Field{
	Type: types.Post,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		return func() (interface{}, error) {
			postId := params.Args["id"].(string)
			if _, err := uuid.Parse(postId); err != nil {
				return nil, errors.New("the given postId is invalid")
			}

			appService := params.Context.Value(config.ServiceKey).(service.Service)
			post, err := appService.GetPostByID(params.Context, postId)
			if err != nil {
				return nil, err
			}

			return post, err
		}, nil
	},
}
