package query

import (
	"net/http"

	"github.com/Ekenzy-101/Go-GraphQL-API/config"
	"github.com/Ekenzy-101/Go-GraphQL-API/graphql/types"
	"github.com/Ekenzy-101/Go-GraphQL-API/service"
	"github.com/graphql-go/graphql"
)

var Me = &graphql.Field{
	Type:        types.User,
	Description: "Returns the authenticated user info if logged in",
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		return func() (interface{}, error) {
			request := params.Context.Value(config.RequestKey).(*http.Request)
			cookie, err := request.Cookie(config.AccessTokenCookieName)
			if err != nil {
				return nil, nil
			}

			appService := params.Context.Value(config.ServiceKey).(service.Service)
			cliams, err := appService.VerifyAccessToken(cookie.Value)
			if err != nil {
				return nil, nil
			}

			user, err := appService.GetUserByID(params.Context, cliams.Subject)
			if err != nil {
				return nil, nil
			}

			return user, nil
		}, nil
	},
}
