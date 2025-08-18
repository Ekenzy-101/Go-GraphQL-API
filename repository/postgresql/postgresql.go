package postgresql

import (
	"context"
	"errors"

	"github.com/Ekenzy-101/Go-GraphQL-API/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4/pgxpool"
)

type postgresRepository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *postgresRepository {
	return &postgresRepository{pool: pool}
}

func (r *postgresRepository) CheckHealth(ctx context.Context) error {
	return r.pool.Ping(ctx)
}

func (r *postgresRepository) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	args := []interface{}{user.Email, user.Name, user.Password}
	dest := []interface{}{&user.ID, &user.CreatedAt}
	sql := `INSERT INTO users (email, name, password) VALUES ($1, $2, $3) RETURNING id, createdAt`
	err := r.pool.QueryRow(ctx, sql, args...).Scan(dest...)
	if isDuplicateKeyError(err) {
		return nil, errors.New("a user with the given email already exists")
	}

	return user, err
}

func (r *postgresRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := &entity.User{Email: email}
	dest := []interface{}{&user.ID, &user.CreatedAt, &user.Name, &user.Password}
	sql := `SELECT id, createdAt, name, password FROM users WHERE email = $1`
	err := r.pool.QueryRow(ctx, sql, email).Scan(dest...)
	if isNotFoundError(err) {
		return nil, errors.New("a user with the given email doesn't exist")
	}

	return user, err
}

func (r *postgresRepository) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("the given userId is invalid")
	}

	user := &entity.User{ID: id}
	dest := []interface{}{&user.CreatedAt, &user.Email, &user.Name, &user.Password}
	sql := `SELECT createdAt, email, name, password FROM users WHERE id = $1`
	err := r.pool.QueryRow(ctx, sql, id).Scan(dest...)
	if isNotFoundError(err) {
		return nil, errors.New("a user with the given id doesn't exist")
	}

	return user, err
}

func (r *postgresRepository) CreatePost(ctx context.Context, post *entity.Post, user *entity.User) (*entity.Post, error) {
	sql := `INSERT INTO posts (content, title, userId) VALUES ($1, $2, $3) RETURNING id, createdAt, updatedAt`
	args := []interface{}{post.Content, post.Title, post.UserID}
	dest := []interface{}{&post.ID, &post.CreatedAt, &post.UpdatedAt}
	return post.SetUser(user), r.pool.QueryRow(ctx, sql, args...).Scan(dest...)
}

func (r *postgresRepository) DeletePostByID(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("the given postId is invalid")
	}

	sql := `DELETE FROM posts WHERE id = $1`
	_, err := r.pool.Exec(ctx, sql, id)
	return err
}

func (r *postgresRepository) GetPostByID(ctx context.Context, id string) (*entity.Post, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("the given postId is invalid")
	}

	post := &entity.Post{ID: id}
	sql := `
	WITH cte_users AS (
		SELECT id, name FROM users
	)
	SELECT p.content, 
	p.createdAt, 
	p.title, 
	p.updatedAt, 
	to_jsonb(u) AS user
	FROM posts AS p
	JOIN cte_users AS u ON p.userId = u.id 
	WHERE p.id = $1`
	dest := []interface{}{&post.Content, &post.CreatedAt, &post.Title, &post.UpdatedAt, &post.User}
	err := r.pool.QueryRow(ctx, sql, id).Scan(dest...)
	if isNotFoundError(err) {
		return nil, errors.New("a post with the given id doesn't exist")
	}

	return post, err
}

func (r *postgresRepository) GetLatestPosts(ctx context.Context, pagination map[string]uint64) ([]entity.Post, error) {
	args := []interface{}{pagination["skip"], pagination["limit"]}
	sql := `
	WITH cte_users AS (
		SELECT id, name FROM users
	)
	SELECT p.id, 
		p.content,
		p.createdAt,
		p.title,
		p.updatedAt, 
		to_jsonb(u) AS user 
	FROM posts AS p 
	JOIN cte_users AS u ON p.userId = u.id 
	ORDER BY p.updatedAt DESC 
	OFFSET $1 ROWS 
	FETCH FIRST $2 ROWS ONLY`
	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	posts := []entity.Post{}
	for rows.Next() {
		post := entity.Post{}
		err := rows.Scan(&post.ID, &post.Content, &post.CreatedAt, &post.Title, &post.UpdatedAt, &post.User)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}
	return posts, nil
}

func (r *postgresRepository) GetUserPosts(ctx context.Context, pagination map[string]uint64, userId string) ([]entity.Post, error) {
	if _, err := uuid.Parse(userId); err != nil {
		return nil, errors.New("the given userId is invalid")
	}

	args := []interface{}{userId, pagination["skip"], pagination["limit"]}
	sql := `
	WITH cte_users AS (
		SELECT id, name FROM users
	)
	SELECT p.id, 
		p.content,
		p.createdAt,
		p.title,
		p.updatedAt, 
		to_jsonb(u) AS user 
	FROM posts AS p 
	JOIN cte_users AS u ON p.userId = u.id
	WHERE p.userId = $1 
	ORDER BY p.updatedAt DESC 
	OFFSET $2 ROWS 
	FETCH FIRST $3 ROWS ONLY`
	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	posts := []entity.Post{}
	for rows.Next() {
		post := entity.Post{}
		err := rows.Scan(&post.ID, &post.Content, &post.CreatedAt, &post.Title, &post.UpdatedAt, &post.User)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}
	return posts, nil
}

func (r *postgresRepository) UpdatePostByID(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	args := []interface{}{post.Content, post.Title, post.UpdatedAt, post.ID}
	sql := `UPDATE posts SET content = $1, title = $2, updatedAt = $3 WHERE id = $4`
	_, err := r.pool.Exec(ctx, sql, args...)
	return post, err
}

func isNotFoundError(err error) bool {
	return err != nil && err.Error() == "no rows in result set"
}

func isDuplicateKeyError(err error) bool {
	pgErr := &pgconn.PgError{}
	return errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation
}
