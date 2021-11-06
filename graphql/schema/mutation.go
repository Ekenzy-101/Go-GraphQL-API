package schema

import (
	"github.com/Ekenzy-101/Go-GraphQL-API/graphql/mutation"
	"github.com/graphql-go/graphql"
)

var Mutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"loginUser":    mutation.LoginUser,
		"logoutUser":   mutation.LogoutUser,
		"registerUser": mutation.RegisterUser,
		"createPost":   mutation.CreatePost,
		"deletePost":   mutation.DeletePost,
		"updatePost":   mutation.UpdatePost,
	},
})
