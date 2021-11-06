package config

type GraphQLContextKey string

const (
	ServiceKey  GraphQLContextKey = "service"
	ResponseKey GraphQLContextKey = "response"
	RequestKey  GraphQLContextKey = "request"
)

const (
	AccessTokenCookieName   = "access_token"
	AccessTokenTTLInSeconds = 60 * 60 * 24
)
