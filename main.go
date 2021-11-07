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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func main() {
	if err := startApplication(); err != nil {
		log.Fatal(err)
	}
}

func startApplication() error {
	ctx := context.Background()
	dbClient, err := setupDatabaseClient(ctx)
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

	repo := repository.New(dbClient)
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

func setupDatabaseClient(ctx context.Context) (interface{}, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.DataBaseURL()))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	collection := client.Database(config.DataBaseName()).Collection(config.UsersCollection)
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "email", Value: bsonx.Int32(1)}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}

	collection = client.Database(config.DataBaseName()).Collection(config.PostsCollection)
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bsonx.Doc{{Key: "userId", Value: bsonx.Int32(1)}},
	})
	if err != nil {
		return nil, err
	}

	return client, nil
	// return pgxpool.Connect(ctx, config.DataBaseURL())
}
