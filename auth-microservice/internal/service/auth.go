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
	ListUsers(ctx context.Context, filter *model.UserFilter) ([]*model.User, int64, error)
	UpdateProfile(ctx context.Context, userID primitive.ObjectID, newName, newEmail string) error
	DeleteProfile(ctx context.Context, userID primitive.ObjectID) error
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

func (s *authService) ListUsers(ctx context.Context, filter *model.UserFilter) ([]*model.User, int64, error) {
	userIDStr, ok := ctx.Value("userID").(string)
	if !ok {
		return nil, 0, errors.New("unauthenticated")
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil, 0, errors.New("invalid user id")
	}

	isAdmin, err := s.IsAdmin(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	if !isAdmin {
		return nil, 0, errors.New("forbidden: only admin can list users")
	}

	users, total, err := repository.ListUsers(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (s *authService) UpdateProfile(ctx context.Context, id primitive.ObjectID, newName, newEmail string) error {

	if newEmail != "" && !utils.ValidEmail(newEmail) {
		return errors.New("invalid email format")
	}

	updateData := make(map[string]interface{})
	if newName == "" {
		return errors.New("Username must not be empty")
	}
	if newEmail == "" {

		return errors.New("Email must not be empty")
	}
	updateData["name"] = newName
	updateData["email"] = newEmail
	if len(updateData) == 0 {
		return errors.New("no data to update")
	}

	return repository.UpdateUser(id, updateData)
}

func (s *authService) DeleteProfile(ctx context.Context, userID primitive.ObjectID) error {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		// ถ้า error มาจาก no documents in result แปลงข้อความ error
		if err.Error() == "mongo: no documents in result" {
			return errors.New("user not found")
		}
		return err
	}
	if user == nil || user.Deleted {
		return errors.New("user not found")
	}

	update := map[string]interface{}{
		"deleted":    true,
		"updated_at": time.Now(),
	}
	return repository.UpdateUser(userID, update)
}
