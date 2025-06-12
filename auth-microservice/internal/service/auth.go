package service

import (
	"auth-microservice/internal/model"
	"auth-microservice/internal/repository"
	"auth-microservice/internal/utils"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService interface {
	Register(ctx context.Context, user *model.User) error
	Login(ctx context.Context, email, password string) (string, error)
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*model.User, error)
	Logout(ctx context.Context, token string) error
}

type authService struct{}

func NewAuthService() AuthService {
	return &authService{}
}

func (s *authService) Register(ctx context.Context, user *model.User) error {
	existing, err := repository.GetUserByEmail(user.Email)
	if err == nil && existing != nil {
		return errors.New("user already exists")
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	user.ID = primitive.NewObjectID()
	user.Deleted = false
	user.CreatedAt = time.Now()

	return repository.CreateUser(user)
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := repository.GetUserByEmail(email)
	if err != nil || user == nil {
		return "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID.Hex())
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) GetUserByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	return repository.GetUserByID(id)
}

func (s *authService) Logout(ctx context.Context, token string) error {
	_, claims, err := utils.ParseJWT(token)
	if err != nil {
		return err
	}

	expUnix := int64(claims["exp"].(float64))
	exp := time.Unix(expUnix, 0)

	return repository.BlacklistToken(token, exp)
}
