package mutation

import (
	"errors"
	"net/http"

	"github.com/Ekenzy-101/Go-GraphQL-API/config"
	"github.com/Ekenzy-101/Go-GraphQL-API/entity"
	"github.com/Ekenzy-101/Go-GraphQL-API/graphql/types"
	"github.com/Ekenzy-101/Go-GraphQL-API/service"
	"github.com/graphql-go/graphql"
)

var CreatePost = &graphql.Field{
	Type: types.Post,
	Args: graphql.FieldConfigArgument{
		"content": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"title": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		request := params.Context.Value(config.RequestKey).(*http.Request)
		cookie, err := request.Cookie(config.AccessTokenCookieName)
		if err != nil {
			return nil, errors.New("please provide a token")
		}

		appService := params.Context.Value(config.ServiceKey).(service.Service)
		cliams, err := appService.VerifyAccessToken(cookie.Value)
		if err != nil {
			return nil, err
		}

		if err := appService.ValidateArgs(params.Context, params.Args); err != nil {
			return nil, err
		}

		post := &entity.Post{
			Content: params.Args["content"].(string),
			Title:   params.Args["title"].(string),
			UserID:  cliams.Subject,
		}
		post, err = appService.CreatePost(params.Context, post, post.UserID)
		if err != nil {
			return nil, err
		}

		return post, nil
	},
}

var DeletePost = &graphql.Field{
	Type: graphql.String,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		request := params.Context.Value(config.RequestKey).(*http.Request)
		cookie, err := request.Cookie(config.AccessTokenCookieName)
		if err != nil {
			return nil, errors.New("please provide a token")
		}

		appService := params.Context.Value(config.ServiceKey).(service.Service)
		cliams, err := appService.VerifyAccessToken(cookie.Value)
		if err != nil {
			return nil, err
		}

		if err := appService.DeletePostIfOwner(params.Context, params.Args["id"].(string), cliams.Subject); err != nil {
			return nil, err
		}

		return "post has been deleted successfully", nil
	},
}

var UpdatePost = &graphql.Field{
	Type: types.Post,
	Args: graphql.FieldConfigArgument{
		"content": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"title": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		request := params.Context.Value(config.RequestKey).(*http.Request)
		cookie, err := request.Cookie(config.AccessTokenCookieName)
		if err != nil {
			return nil, errors.New("please provide a token")
		}

		appService := params.Context.Value(config.ServiceKey).(service.Service)
		cliams, err := appService.VerifyAccessToken(cookie.Value)
		if err != nil {
			return nil, err
		}

		if err := appService.ValidateArgs(params.Context, params.Args); err != nil {
			return nil, err
		}

		post, err := appService.UpdatePostIfOwner(params.Context, params.Args, cliams.Subject)
		if err != nil {
			return nil, err
		}

		return post, nil
	},
}
