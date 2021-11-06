package mutation

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Ekenzy-101/Go-GraphQL-API/config"
	"github.com/Ekenzy-101/Go-GraphQL-API/entity"
	"github.com/Ekenzy-101/Go-GraphQL-API/graphql/types"
	"github.com/Ekenzy-101/Go-GraphQL-API/service"
	"github.com/graphql-go/graphql"
)

var LoginUser = &graphql.Field{
	Type: types.User,
	Args: graphql.FieldConfigArgument{
		"email": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"password": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		appService := params.Context.Value(config.ServiceKey).(service.Service)
		if err := appService.ValidateArgs(params.Context, params.Args); err != nil {
			return nil, err
		}

		user, err := appService.GetUserByEmail(params.Context, params.Args["email"].(string))
		if err != nil && strings.Contains(err.Error(), "email") {
			return nil, errors.New("invalid email or password")
		}

		if err != nil {
			return nil, err
		}

		matches, err := appService.ComparePassword(params.Args["password"].(string), user.Password)
		if err != nil {
			return nil, err
		}

		if !matches {
			return nil, errors.New("invalid email or password")
		}

		accessToken, err := appService.GenerateAccessToken(user)
		if err != nil {
			return nil, err
		}

		response := params.Context.Value(config.ResponseKey).(http.ResponseWriter)
		http.SetCookie(response, &http.Cookie{
			HttpOnly: true,
			MaxAge:   config.AccessTokenTTLInSeconds,
			Name:     config.AccessTokenCookieName,
			Path:     "/",
			Secure:   config.IsProduction(),
			Value:    accessToken,
		})
		return user.SetPassword(""), nil
	},
}

var LogoutUser = &graphql.Field{
	Type: graphql.String,
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		response := params.Context.Value(config.ResponseKey).(http.ResponseWriter)
		http.SetCookie(response, &http.Cookie{
			HttpOnly: true,
			MaxAge:   -1,
			Name:     config.AccessTokenCookieName,
			Path:     "/",
			Secure:   config.IsProduction(),
			Value:    "",
		})
		return "logout successfully", nil
	},
}

var RegisterUser = &graphql.Field{
	Type: types.User,
	Args: graphql.FieldConfigArgument{
		"email": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"name": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"password": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		appService := params.Context.Value(config.ServiceKey).(service.Service)
		err := appService.ValidateArgs(params.Context, params.Args)
		if err != nil {
			return nil, err
		}

		user := &entity.User{
			Email:    params.Args["email"].(string),
			Name:     params.Args["name"].(string),
			Password: params.Args["password"].(string),
		}
		user.Password, err = appService.HashPassword(user.Password)
		if err != nil {
			return nil, err
		}

		user, err = appService.CreateUser(params.Context, user)
		if err != nil {
			return nil, err
		}

		accessToken, err := appService.GenerateAccessToken(user)
		if err != nil {
			return nil, err
		}

		response := params.Context.Value(config.ResponseKey).(http.ResponseWriter)
		http.SetCookie(response, &http.Cookie{
			HttpOnly: true,
			MaxAge:   config.AccessTokenTTLInSeconds,
			Name:     config.AccessTokenCookieName,
			Path:     "/",
			Secure:   config.IsProduction(),
			Value:    accessToken,
		})
		return user.SetPassword(""), nil
	},
}
