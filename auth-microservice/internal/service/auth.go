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

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	activeToken, err := repository.GetActiveTokenByUserID(user.ID)
	if err == nil && activeToken != "" {
		// เช็คว่า token ยังไม่หมดอายุ และยังไม่ถูก blacklist
		if !utils.IsTokenExpired(activeToken) && !repository.IsTokenBlacklisted(activeToken) {
			// ถ้ายัง valid ให้ return token เดิม
			return activeToken, nil
		}
		// ถ้า token หมดอายุ หรือ ถูก blacklist ก็ไปสร้างใหม่
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
