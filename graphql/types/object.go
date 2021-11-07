package types

import (
	"github.com/Ekenzy-101/Go-GraphQL-API/config"
	"github.com/Ekenzy-101/Go-GraphQL-API/entity"
	"github.com/Ekenzy-101/Go-GraphQL-API/service"
	"github.com/graphql-go/graphql"
)

var PaginationArgs = graphql.FieldConfigArgument{
	"skip": &graphql.ArgumentConfig{
		Type: UInt,
	},
	"limit": &graphql.ArgumentConfig{
		Type: UInt,
	},
}

var Post = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Post",
	Description: "A post type",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Description: "Unique identifier of the post",
			Type:        graphql.NewNonNull(graphql.ID),
		},
		"createdAt": &graphql.Field{
			Description: "Time when the post was created",
			Type:        graphql.NewNonNull(graphql.DateTime),
		},
		"content": &graphql.Field{
			Description: "Content of the post",
			Type:        graphql.NewNonNull(graphql.String),
		},
		"title": &graphql.Field{
			Description: "Title of the post",
			Type:        graphql.NewNonNull(graphql.String),
		},
		"updatedAt": &graphql.Field{
			Description: "Time when the post was updated",
			Type:        graphql.NewNonNull(graphql.DateTime),
		},
	},
})

var User = graphql.NewObject(graphql.ObjectConfig{
	Name:        "User",
	Description: "A user type",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Description: "Unique identifier of the user",
			Type:        graphql.NewNonNull(graphql.ID),
		},
		"createdAt": &graphql.Field{
			Description: "Time when the user was registered",
			Type:        graphql.NewNonNull(graphql.DateTime),
		},
		"email": &graphql.Field{
			Description: "Email address of the user",
			Type:        graphql.NewNonNull(graphql.String),
		},
		"name": &graphql.Field{
			Description: "Fullname of the user",
			Type:        graphql.NewNonNull(graphql.String),
		},
		"posts": &graphql.Field{
			Description: "Posts of the user",
			Type:        graphql.NewList(Post),
			Args:        PaginationArgs,
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return func() (interface{}, error) {
					user := params.Source.(*entity.User)
					if user == nil {
						return nil, nil
					}

					skip, ok := params.Args["skip"].(uint64)
					if !ok {
						skip = 0
					}

					limit, ok := params.Args["limit"].(uint64)
					if !ok {
						limit = 2
					}

					if len(params.Info.Path.AsArray()) > 2 {
						return []interface{}{}, nil
					}

					appService := params.Context.Value(config.ServiceKey).(service.Service)
					posts, err := appService.GetUserPosts(params.Context, map[string]uint64{"skip": skip, "limit": limit}, user.ID)
					if err != nil {
						return nil, err
					}

					return posts, nil
				}, nil
			},
		},
	},
})

func init() {
	Post.AddFieldConfig("user", &graphql.Field{
		Description: "User of the post",
		Type:        graphql.NewNonNull(User),
	})
}
