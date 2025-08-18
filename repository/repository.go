package repository

import (
	"context"
	"errors"

	"github.com/Ekenzy-101/Go-GraphQL-API/config"
	"github.com/Ekenzy-101/Go-GraphQL-API/entity"
	"github.com/Ekenzy-101/Go-GraphQL-API/repository/mongodb"
	"github.com/Ekenzy-101/Go-GraphQL-API/repository/postgresql"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
}

type PostRepository interface {
	CreatePost(ctx context.Context, post *entity.Post, user *entity.User) (*entity.Post, error)
	DeletePostByID(ctx context.Context, id string) error
	GetLatestPosts(ctx context.Context, pagination map[string]uint64) ([]entity.Post, error)
	GetUserPosts(ctx context.Context, pagination map[string]uint64, userId string) ([]entity.Post, error)
	GetPostByID(ctx context.Context, id string) (*entity.Post, error)
	UpdatePostByID(ctx context.Context, post *entity.Post) (*entity.Post, error)
}

type Repository interface {
	PostRepository
	UserRepository
	CheckHealth(ctx context.Context) error
}

func New(ctx context.Context) (Repository, error) {
	cacheClient, err := setupCacheClient(ctx)
	if err != nil {
		return nil, err
	}

	switch config.DataBaseType() {
	case config.DataBaseMongo:
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

		return mongodb.New(client, cacheClient), nil
	case config.DataBasePostgres:
		client, err := pgxpool.Connect(ctx, config.DataBaseURL())
		if err != nil {
			return nil, err
		}

		return postgresql.New(client), nil
	default:
		return nil, errors.New("unknown database type")
	}
}

func setupCacheClient(ctx context.Context) (*redis.Client, error) {
	options, err := redis.ParseURL(config.CacheURL())
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
