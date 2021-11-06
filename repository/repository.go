package repository

import (
	"context"

	"github.com/Ekenzy-101/Go-GraphQL-API/entity"
	"github.com/Ekenzy-101/Go-GraphQL-API/repository/postgresql"
	"github.com/jackc/pgx/v4/pgxpool"
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
}

func New(dbClient interface{}) Repository {
	return postgresql.New(dbClient.(*pgxpool.Pool))
}
