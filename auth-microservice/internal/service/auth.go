package service

import (
	"auth-microservice/internal/model"
	"auth-microservice/internal/repository"
	"auth-microservice/internal/utils"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService interface {
	Register(ctx context.Context, user *model.User) error
	Login(ctx context.Context, email, password string) (string, error)
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*model.User, error)
	Logout(ctx context.Context, token string) error
	IsAdmin(ctx context.Context, userID primitive.ObjectID) (bool, error)
	AddRole(ctx context.Context, adminUserID, targetUserID primitive.ObjectID, newRole string) error
}

type authService struct{}

func NewAuthService() AuthService {
	return &authService{}
}

func (s *authService) Register(ctx context.Context, user *model.User) error {
	if user.Email == "" || user.Password == "" || user.Name == "" {
		return errors.New("email, password and name must not be empty")
	}

	if !utils.ValidEmail(user.Email) {
		return errors.New("invalid email format")
	}

	if !utils.ValidPassword(user.Password) {
		return errors.New("password must be at least 8 characters long and include uppercase, lowercase, and number")
	}

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
	if email == "" || password == "" {
		return "", errors.New("email and password must not be empty")
	}

	user, err := repository.GetUserByEmail(email)
	if err != nil || user == nil {
		return "", errors.New("invalid credentials")
	}

	key := fmt.Sprintf("login_attempts:%s", email)
	allowed, err := utils.RateLimit(ctx, key, 5, time.Minute)
	if err != nil {
		return "", err
	}
	if !allowed {
		return "", errors.New("too many login attempts, please try again later")
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

func (s *authService) IsAdmin(ctx context.Context, id primitive.ObjectID) (bool, error) {
	user, err := repository.GetUserByID(id)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, errors.New("user not found")
	}
	return user.Role == "admin", nil
}

func (s *authService) AddRole(ctx context.Context, adminUserID, targetUserID primitive.ObjectID, newRole string) error {
	isAdmin, err := s.IsAdmin(ctx, adminUserID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return errors.New("forbidden: only admin can update role")
	}

	if newRole != "user" && newRole != "admin" {
		return errors.New("invalid role: must be 'user' or 'admin'")
	}

	updateData := map[string]interface{}{
		"role": newRole,
	}

	return repository.UpdateUser(targetUserID, updateData)
}
