package repository

import (
	"auth-microservice/internal/db"
	"auth-microservice/internal/model"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateUser(user *model.User) error {
	collection := db.GetCollection(db.DB_NAME, db.USER_COLLECTION)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, user)
	return err
}

func GetUserByEmail(email string) (*model.User, error) {
	collection := db.GetCollection(db.DB_NAME, db.USER_COLLECTION)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user model.User
	err := collection.FindOne(ctx, bson.M{"email": email, "deleted": false}).Decode(&user)
	return &user, err
}

func GetUserByID(id primitive.ObjectID) (*model.User, error) {
	collection := db.GetCollection(db.DB_NAME, db.USER_COLLECTION)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user model.User
	err := collection.FindOne(ctx, bson.M{"_id": id, "deleted": false}).Decode(&user)
	return &user, err
}

func UpdateUser(userID primitive.ObjectID, updateData map[string]interface{}) error {
	collection := db.GetCollection(db.DB_NAME, db.USER_COLLECTION)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": userID, "deleted": false}
	update := bson.M{"$set": updateData}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func ListUsers(ctx context.Context, filter *model.UserFilter) ([]*model.User, int64, error) {
	collection := db.GetUserCollection()

	bsonFilter := bson.M{}
	if filter.Name != "" {
		bsonFilter["name"] = filter.Name
	}
	if filter.Email != "" {
		bsonFilter["email"] = filter.Email
	}
	bsonFilter["deleted"] = false

	limit := filter.Limit
	if limit <= 0 {
		limit = 10
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	skip := (page - 1) * limit

	total, err := collection.CountDocuments(ctx, bsonFilter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find()
	opts.SetSkip(skip)
	opts.SetLimit(limit)

	cursor, err := collection.Find(ctx, bsonFilter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []*model.User
	for cursor.Next(ctx) {
		var u model.User
		if err := cursor.Decode(&u); err != nil {
			log.Printf("failed to decode user: %v", err)
			continue
		}
		users = append(users, &u)
	}

	return users, total, nil
}
