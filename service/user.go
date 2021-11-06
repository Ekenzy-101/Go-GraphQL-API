package service

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Ekenzy-101/Go-GraphQL-API/config"
	"github.com/Ekenzy-101/Go-GraphQL-API/entity"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/argon2"
)

type UserRepository interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
	GetByID(ctx context.Context, id string) (entity.User, error)
}

type passwordConfig struct {
	time      uint32
	memory    uint32
	threads   uint8
	keyLength uint32
}

type jwtOptions struct {
	signingMethod jwt.SigningMethod
	claims        jwt.Claims
	secret        string
	token         string
}

func (s *service) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	return s.repo.CreateUser(ctx, user)
}

func (s *service) ComparePassword(password string, hashedPassword string) (bool, error) {
	parts := strings.Split(hashedPassword, "$")
	if len(parts) < 4 {
		return false, errors.New("invalid password hash")
	}

	c := &passwordConfig{}
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &c.memory, &c.time, &c.threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	c.keyLength = uint32(len(decodedHash))
	comparisonHash := argon2.IDKey([]byte(password), salt, c.time, c.memory, c.threads, c.keyLength)
	return (subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1), nil
}

func (s *service) HashPassword(password string) (string, error) {
	c := &passwordConfig{
		time:      1,
		memory:    64 * 1024,
		threads:   4,
		keyLength: 32,
	}
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, c.time, c.memory, c.threads, c.keyLength)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	return fmt.Sprintf(format, argon2.Version, c.memory, c.time, c.threads, b64Salt, b64Hash), nil
}

func (s *service) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *service) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *service) GenerateAccessToken(user *entity.User) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Second * config.AccessTokenTTLInSeconds)),
		Subject:   user.ID,
	}

	return signJWTToken(jwtOptions{
		signingMethod: jwt.SigningMethodHS256,
		claims:        claims,
		secret:        config.AccessTokenSecret(),
	})
}

func (s *service) VerifyAccessToken(token string) (*jwt.RegisteredClaims, error) {
	cliams, err := verifyJWTToken(jwtOptions{
		claims: &jwt.RegisteredClaims{},
		secret: config.AccessTokenSecret(),
		token:  token,
	})
	if err != nil {
		return nil, err
	}

	accessTokenClaims, ok := cliams.(*jwt.RegisteredClaims)
	if !ok {
		return nil, errors.New("claims' format is invalid")
	}

	return accessTokenClaims, nil
}

func signJWTToken(options jwtOptions) (string, error) {
	token := jwt.NewWithClaims(options.signingMethod, options.claims)
	signedToken, err := token.SignedString([]byte(options.secret))

	return signedToken, err
}

func verifyJWTToken(options jwtOptions) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(
		options.token,
		options.claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(options.secret), nil
		},
	)

	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	return token.Claims, nil
}
