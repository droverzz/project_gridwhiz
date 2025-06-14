package service

import (
	"auth-microservice/internal/db"
	"auth-microservice/internal/model"
	"auth-microservice/internal/repository"
	"auth-microservice/internal/utils"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrForbidden          = errors.New("forbidden")
	ErrUnauthenticated    = errors.New("unauthenticated")
	ErrNotFound           = errors.New("not found")
	ErrInvalidArgument    = errors.New("invalid argument")
)

type AuthService interface {
	Register(ctx context.Context, user *model.User) error
	Login(ctx context.Context, email, password string) (string, error)
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*model.User, error)
	Logout(ctx context.Context, token string) error
	IsAdmin(ctx context.Context, userID primitive.ObjectID) (bool, error)
	AddRole(ctx context.Context, adminUserID, targetUserID primitive.ObjectID, newRole string) error
	ListUsers(ctx context.Context, filter *model.UserFilter) ([]*model.User, int64, error)
	UpdateProfile(ctx context.Context, userID primitive.ObjectID, newName, newEmail string) error
	DeleteProfile(ctx context.Context, userID primitive.ObjectID) error
	GeneratePasswordResetToken(ctx context.Context, userID primitive.ObjectID) (string, error)
	ResetPassword(ctx context.Context, resetToken, newPassword string) error
}

type authService struct{}

func NewAuthService() AuthService {
	return &authService{}
}

func (s *authService) Register(ctx context.Context, user *model.User) error {
	if user.Email == "" || user.Password == "" || user.Name == "" {
		return ErrInvalidArgument
	}

	if !utils.ValidEmail(user.Email) {
		return ErrInvalidArgument
	}

	if !utils.ValidPassword(user.Password) {
		return ErrInvalidArgument
	}

	existing, err := repository.GetUserByEmail(user.Email)
	if err == nil && existing != nil {
		return ErrUserExists
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
		return "", ErrInvalidArgument
	}

	user, err := repository.GetUserByEmail(email)
	if err != nil || user == nil {
		return "", ErrInvalidCredentials
	}

	key := fmt.Sprintf("login_attempts:%s", email)
	allowed, err := utils.RateLimit(ctx, key, 5, time.Minute)
	if err != nil {
		return "", err
	}
	if !allowed {
		return "", ErrForbidden
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", ErrInvalidCredentials
	}

	token, err := utils.GenerateJWT(user.ID.Hex())
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) GetUserByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	user, err := repository.GetUserByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	return user, nil
}

func (s *authService) Logout(ctx context.Context, token string) error {
	_, claims, err := utils.ParseJWT(token)
	if err != nil {
		return ErrUnauthenticated
	}

	expUnix := int64(claims["exp"].(float64))
	exp := time.Unix(expUnix, 0)

	return repository.BlacklistToken(token, exp)
}

func (s *authService) IsAdmin(ctx context.Context, id primitive.ObjectID) (bool, error) {
	user, err := repository.GetUserByID(id)
	if err != nil {
		return false, ErrNotFound
	}
	if user == nil {
		return false, ErrNotFound
	}
	return user.Role == "admin", nil
}

func (s *authService) AddRole(ctx context.Context, adminUserID, targetUserID primitive.ObjectID, newRole string) error {
	isAdmin, err := s.IsAdmin(ctx, adminUserID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return ErrForbidden
	}

	if newRole != "user" && newRole != "admin" {
		return ErrInvalidArgument
	}

	updateData := map[string]interface{}{
		"role": newRole,
	}

	return repository.UpdateUser(targetUserID, updateData)
}

func (s *authService) ListUsers(ctx context.Context, filter *model.UserFilter) ([]*model.User, int64, error) {
	userIDStr, ok := ctx.Value("user_id").(string)
	if !ok {
		return nil, 0, ErrUnauthenticated
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil, 0, ErrInvalidArgument
	}

	isAdmin, err := s.IsAdmin(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	if !isAdmin {
		return nil, 0, ErrForbidden
	}

	users, total, err := repository.ListUsers(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (s *authService) UpdateProfile(ctx context.Context, id primitive.ObjectID, newName, newEmail string) error {
	if newEmail != "" && !utils.ValidEmail(newEmail) {
		return ErrInvalidArgument
	}

	updateData := make(map[string]interface{})
	if newName == "" {
		return ErrInvalidArgument
	}
	if newEmail == "" {
		return ErrInvalidArgument
	}
	updateData["name"] = newName
	updateData["email"] = newEmail
	if len(updateData) == 0 {
		return ErrInvalidArgument
	}

	return repository.UpdateUser(id, updateData)
}

func (s *authService) DeleteProfile(ctx context.Context, userID primitive.ObjectID) error {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return ErrNotFound
		}
		return err
	}
	if user == nil || user.Deleted {
		return ErrNotFound
	}

	update := map[string]interface{}{
		"deleted":    true,
		"updated_at": time.Now(),
	}
	return repository.UpdateUser(userID, update)
}

func (s *authService) GeneratePasswordResetToken(ctx context.Context, userID primitive.ObjectID) (string, error) {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return "", err
	}
	email := user.Email
	fmt.Println(email)

	token, err := utils.GenerateResetToken(user.ID.Hex())
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) ResetPassword(ctx context.Context, resetToken, newPassword string) error {
	userIDHex, err := utils.ValidateResetToken(resetToken)
	if err != nil {
		return ErrInvalidArgument
	}

	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return ErrInvalidArgument
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	err = db.UpdatePassword(ctx, userID, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}
