package db

import (
	"auth-microservice/internal/model"
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

var (
	DB_NAME         string
	USER_COLLECTION string
)

func InitMongoDB(uri string) error {
	DB_NAME = os.Getenv("DB_NAME")
	USER_COLLECTION = os.Getenv("USER_COLLECTION")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	Client = client
	log.Println("Connected to MongoDB!")
	return nil
}

func GetCollection(database, collection string) *mongo.Collection {
	return Client.Database(database).Collection(collection)
}

func GetUserCollection() *mongo.Collection {
	return GetCollection(DB_NAME, USER_COLLECTION)
}
func FindUsers(ctx context.Context, filter model.UserFilter) ([]model.User, int64, error) {
	collection := GetUserCollection()

	bsonFilter := bson.M{}
	if filter.Name != "" {
		bsonFilter["name"] = bson.M{"$regex": filter.Name, "$options": "i"}
	}
	if filter.Email != "" {
		bsonFilter["email"] = bson.M{"$regex": filter.Email, "$options": "i"}
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 10
	}

	page := filter.Page
	if page <= 0 {
		page = 1
	}
	skip := (page - 1) * limit

	// นับจำนวนเอกสารทั้งหมดที่ตรงเงื่อนไข
	total, err := collection.CountDocuments(ctx, bsonFilter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().SetLimit(limit).SetSkip(skip)

	cursor, err := collection.Find(ctx, bsonFilter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []model.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
