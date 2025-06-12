package repository

import (
	"auth-microservice/internal/db"
	models "auth-microservice/internal/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	DB_NAME         = "gridwhizdb"
	USER_COLLECTION = "users"
)

func CreateUser(user *models.User) error {
	collection := db.GetCollection(DB_NAME, USER_COLLECTION)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, user)
	return err
}

func GetUserByEmail(email string) (*models.User, error) {
	collection := db.GetCollection(DB_NAME, USER_COLLECTION)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email, "deleted": false}).Decode(&user)
	return &user, err
}

func GetUserByID(id primitive.ObjectID) (*models.User, error) {
	collection := db.GetCollection(DB_NAME, USER_COLLECTION)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"_id": id, "deleted": false}).Decode(&user)
	return &user, err
}
