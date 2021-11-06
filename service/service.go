package service

import (
	"context"

	"github.com/Ekenzy-101/Go-GraphQL-API/entity"
	"github.com/Ekenzy-101/Go-GraphQL-API/repository"
	"github.com/golang-jwt/jwt/v4"
)

type UserService interface {
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	ComparePassword(password string, hashedPassword string) (bool, error)
	GenerateAccessToken(user *entity.User) (string, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
	HashPassword(password string) (string, error)
	VerifyAccessToken(token string) (*jwt.RegisteredClaims, error)
}

type PostService interface {
	CreatePost(ctx context.Context, post *entity.Post, userId string) (*entity.Post, error)
	DeletePostIfOwner(ctx context.Context, postId, userId string) error
	GetLatestPosts(ctx context.Context, pagination map[string]uint64) ([]entity.Post, error)
	GetUserPosts(ctx context.Context, pagination map[string]uint64, userId string) ([]entity.Post, error)
	GetPostByID(ctx context.Context, id string) (*entity.Post, error)
	UpdatePostIfOwner(ctx context.Context, args map[string]interface{}, userId string) (*entity.Post, error)
}

type CommonService interface {
	ValidateArgs(ctx context.Context, args map[string]interface{}) error
}

type Service interface {
	CommonService
	PostService
	UserService
}

type service struct {
	repo repository.Repository
}

func New(repo repository.Repository) Service {
	return &service{repo: repo}
}
