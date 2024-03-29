package mongodb

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Ekenzy-101/Go-GraphQL-API/config"
	"github.com/Ekenzy-101/Go-GraphQL-API/entity"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepository struct {
	dbClient    *mongo.Client
	cacheClient *redis.Client
}

const (
	LatestPostsKey = "latestposts"
)

func New(dbClient *mongo.Client, cacheClient *redis.Client) *mongoRepository {
	return &mongoRepository{dbClient: dbClient, cacheClient: cacheClient}
}

func (r *mongoRepository) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	user.SetID(primitive.NewObjectID().Hex()).SetCreatedAt(time.Now().UTC())

	collection := r.dbClient.Database(config.DataBaseName()).Collection(config.UsersCollection)
	_, err := collection.InsertOne(ctx, user)
	if isDuplicateKeyError(err) {
		return nil, errors.New("a user with the given email already exists")
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *mongoRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := &entity.User{}

	collection := r.dbClient.Database(config.DataBaseName()).Collection(config.UsersCollection)
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(user)
	if isNotFoundError(err) {
		return nil, errors.New("a user with the given email doesn't exist")
	}

	return user, err
}

func (r *mongoRepository) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	user := &entity.User{}

	collection := r.dbClient.Database(config.DataBaseName()).Collection(config.UsersCollection)
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(user)
	if isNotFoundError(err) {
		return nil, errors.New("a user with the given id doesn't exist")
	}

	return user, err
}

func (r *mongoRepository) CreatePost(ctx context.Context, post *entity.Post, user *entity.User) (*entity.Post, error) {
	now := time.Now().UTC()
	post.SetID(primitive.NewObjectID().Hex()).SetCreatedAt(now).SetUpdatedAt(now)

	collection := r.dbClient.Database(config.DataBaseName()).Collection(config.PostsCollection)
	_, err := collection.InsertOne(ctx, post)
	if err != nil {
		return nil, err
	}

	return post.SetUser(user), r.deleteLatestPostsFromCache(ctx)
}

func (r *mongoRepository) DeletePostByID(ctx context.Context, id string) error {
	collection := r.dbClient.Database(config.DataBaseName()).Collection(config.PostsCollection)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *mongoRepository) GetPostByID(ctx context.Context, id string) (*entity.Post, error) {
	pipeline := bson.A{
		bson.M{"$match": bson.M{"_id": id}},
		bson.M{
			"$lookup": bson.M{
				"from": config.UsersCollection,
				"let":  bson.M{"userId": "$userId"},
				"pipeline": bson.A{
					bson.M{
						"$match": bson.M{
							"$expr": bson.M{"$eq": bson.A{"$_id", "$$userId"}}},
					},
					bson.M{
						"$project": bson.M{"name": 1},
					},
				},
				"as": "user",
			},
		},
		bson.M{"$unwind": "$user"},
		bson.M{"$project": bson.M{"userId": 0}},
	}

	collection := r.dbClient.Database(config.DataBaseName()).Collection(config.PostsCollection)
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	posts := make([]entity.Post, 0, 1)
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, err
	}

	if len(posts) == 0 {
		return nil, errors.New("a post with the given id doesn't exist")
	}

	return &posts[0], err
}

func (r *mongoRepository) GetLatestPosts(ctx context.Context, pagination map[string]uint64) ([]entity.Post, error) {
	cachedPosts, err := r.getLatestPostsFromCache(ctx, pagination)
	if !isNotFoundError(err) && err != nil {
		return nil, err
	}

	if err == nil {
		return cachedPosts, nil
	}

	pipeline := bson.A{
		bson.M{"$sort": bson.M{"updatedAt": -1}},
		bson.M{"$skip": pagination["skip"]},
		bson.M{"$limit": pagination["limit"]},
		bson.M{
			"$lookup": bson.M{
				"from": config.UsersCollection,
				"let":  bson.M{"userId": "$userId"},
				"pipeline": bson.A{
					bson.M{
						"$match": bson.M{
							"$expr": bson.M{"$eq": bson.A{"$_id", "$$userId"}}},
					},
					bson.M{
						"$project": bson.M{"name": 1},
					},
				},
				"as": "user",
			},
		},
		bson.M{"$unwind": "$user"},
		bson.M{"$project": bson.M{"userId": 0}},
	}
	collection := r.dbClient.Database(config.DataBaseName()).Collection(config.PostsCollection)
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	posts := []entity.Post{}
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, err
	}

	return posts, r.saveLatestPostsToCache(ctx, pagination, posts)
}

func (r *mongoRepository) GetUserPosts(ctx context.Context, pagination map[string]uint64, userId string) ([]entity.Post, error) {
	pipeline := bson.A{
		bson.M{"$match": bson.M{"userId": userId}},
		bson.M{"$sort": bson.M{"updatedAt": -1}},
		bson.M{"$skip": pagination["skip"]},
		bson.M{"$limit": pagination["limit"]},
		bson.M{
			"$lookup": bson.M{
				"from": config.UsersCollection,
				"let":  bson.M{"userId": "$userId"},
				"pipeline": bson.A{
					bson.M{
						"$match": bson.M{
							"$expr": bson.M{"$eq": bson.A{"$_id", "$$userId"}}},
					},
					bson.M{
						"$project": bson.M{"name": 1},
					},
				},
				"as": "user",
			},
		},
		bson.M{"$unwind": "$user"},
		bson.M{"$project": bson.M{"userId": 0}},
	}
	collection := r.dbClient.Database(config.DataBaseName()).Collection(config.PostsCollection)
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	posts := []entity.Post{}
	return posts, cursor.All(ctx, &posts)
}

func (r *mongoRepository) UpdatePostByID(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	update := bson.M{
		"$set": bson.M{
			"content":   post.Content,
			"title":     post.Title,
			"updatedAt": post.UpdatedAt,
		},
	}
	collection := r.dbClient.Database(config.DataBaseName()).Collection(config.PostsCollection)
	_, err := collection.UpdateByID(ctx, post.ID, update)
	return post, err
}

func (r *mongoRepository) getLatestPostsFromCache(ctx context.Context, pagination map[string]uint64) ([]entity.Post, error) {
	field, err := json.Marshal(pagination)
	if err != nil {
		return nil, err
	}

	data, err := r.cacheClient.HGet(ctx, LatestPostsKey, string(field)).Bytes()
	if err != nil {
		return nil, err
	}

	posts := make([]entity.Post, 0, pagination["limit"])
	if err := json.Unmarshal(data, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *mongoRepository) deleteLatestPostsFromCache(ctx context.Context) error {
	return r.cacheClient.Del(ctx, LatestPostsKey).Err()
}

func (r *mongoRepository) saveLatestPostsToCache(ctx context.Context, pagination map[string]uint64, posts []entity.Post) error {
	key, err := json.Marshal(pagination)
	if err != nil {
		return err
	}

	value, err := json.Marshal(posts)
	if err != nil {
		return err
	}

	if err := r.cacheClient.HSet(ctx, LatestPostsKey, string(key), string(value)).Err(); err != nil {
		return err
	}

	return nil
}

func isNotFoundError(err error) bool {
	return errors.Is(err, mongo.ErrNoDocuments) || errors.Is(err, redis.Nil)
}

func isDuplicateKeyError(err error) bool {
	return mongo.IsDuplicateKeyError(err)
}
