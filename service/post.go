package service

import (
	"context"
	"errors"
	"time"

	"github.com/Ekenzy-101/Go-GraphQL-API/entity"
)

func (s *service) CreatePost(ctx context.Context, post *entity.Post, userId string) (*entity.Post, error) {
	user, err := s.repo.GetUserByID(ctx, userId)
	if err != nil {
		return nil, err
	}
	return s.repo.CreatePost(ctx, post, user)
}

func (s *service) DeletePostIfOwner(ctx context.Context, postId, userId string) error {
	post, err := s.repo.GetPostByID(ctx, postId)
	if err != nil {
		return err
	}

	if post.User == nil || post.User.ID != userId {
		return errors.New("you are not allowed to delete this post")
	}

	return s.repo.DeletePostByID(ctx, postId)
}

func (s *service) GetPostByID(ctx context.Context, id string) (*entity.Post, error) {
	return s.repo.GetPostByID(ctx, id)
}

func (s *service) GetLatestPosts(ctx context.Context, pagination map[string]uint64) ([]entity.Post, error) {
	return s.repo.GetLatestPosts(ctx, pagination)
}

func (s *service) GetUserPosts(ctx context.Context, pagination map[string]uint64, userId string) ([]entity.Post, error) {
	return s.repo.GetUserPosts(ctx, pagination, userId)
}

func (s *service) UpdatePostIfOwner(ctx context.Context, args map[string]interface{}, userId string) (*entity.Post, error) {
	post, err := s.repo.GetPostByID(ctx, args["id"].(string))
	if err != nil {
		return nil, err
	}

	if post.User == nil || post.User.ID != userId {
		return nil, errors.New("you are not allowed to update this post")
	}

	post.SetContent(args["content"].(string)).SetTitle(args["title"].(string)).SetUpdatedAt(time.Now().UTC())
	return s.repo.UpdatePostByID(ctx, post)
}
