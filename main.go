package main

import (
	"context"
	"log"

	"github.com/Ekenzy-101/Go-GraphQL-API/config"
	"github.com/Ekenzy-101/Go-GraphQL-API/graphql/schema"
	"github.com/Ekenzy-101/Go-GraphQL-API/repository"
	"github.com/Ekenzy-101/Go-GraphQL-API/service"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	if err := startApplication(); err != nil {
		log.Fatal(err)
	}
}

func startApplication() error {
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, config.DataBaseURL())
	if err != nil {
		return err
	}

	graphqlSchema, err := schema.New()
	if err != nil {
		return err
	}

	graphqlHandler := handler.New(&handler.Config{
		Schema:     &graphqlSchema,
		Pretty:     true,
		Playground: true,
	})

	repo := repository.New(pool)
	appService := service.New(repo)
	router := gin.Default()
	router.POST("/graphql", func(c *gin.Context) {
		ctx := context.WithValue(context.Background(), config.ServiceKey, appService)
		ctx = context.WithValue(ctx, config.ResponseKey, c.Writer)
		ctx = context.WithValue(ctx, config.RequestKey, c.Request)

		graphqlHandler.ContextHandler(ctx, c.Writer, c.Request)
	})
	return router.Run(":5000")
}
