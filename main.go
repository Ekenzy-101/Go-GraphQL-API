package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Ekenzy-101/Go-GraphQL-API/config"
	"github.com/Ekenzy-101/Go-GraphQL-API/graphql/schema"
	"github.com/Ekenzy-101/Go-GraphQL-API/repository"
	"github.com/Ekenzy-101/Go-GraphQL-API/service"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
)

func main() {
	if err := startApplication(); err != nil {
		log.Fatal(err)
	}
}

func startApplication() error {
	ctx := context.Background()

	graphqlSchema, err := schema.New()
	if err != nil {
		return err
	}

	graphqlHandler := handler.New(&handler.Config{
		Schema:     &graphqlSchema,
		Pretty:     true,
		Playground: true,
	})

	repo, err := repository.New(ctx)
	if err != nil {
		return err
	}

	appService := service.New(repo)

	router := gin.Default()
	router.GET("/healthcheck", func(c *gin.Context) {
		if err := repo.CheckHealth(c.Request.Context()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	})
	router.Any("/graphql", func(c *gin.Context) {
		ctx := context.WithValue(context.Background(), config.ServiceKey, appService)
		ctx = context.WithValue(ctx, config.ResponseKey, c.Writer)
		ctx = context.WithValue(ctx, config.RequestKey, c.Request)

		graphqlHandler.ContextHandler(ctx, c.Writer, c.Request)
	})
	return router.Run(":" + config.Port())
}
