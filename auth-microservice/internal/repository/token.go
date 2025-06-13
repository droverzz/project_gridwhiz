package repository

import (
	"auth-microservice/internal/db"
	"auth-microservice/internal/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

const BLACKLIST_COLLECTION = "blacklisted_tokens"

func BlacklistToken(token string, exp time.Time) error {
	col := db.GetCollection(db.DB_NAME, BLACKLIST_COLLECTION)

	count, err := col.CountDocuments(context.TODO(), bson.M{"token": token})
	if err != nil {
		return err
	}
	if count > 0 {

		return nil
	}

	_, err = col.InsertOne(context.TODO(), model.BlacklistedToken{
		Token:     token,
		ExpiredAt: exp,
	})
	return err
}

func IsTokenBlacklisted(token string) (bool, error) {
	col := db.GetCollection(db.DB_NAME, BLACKLIST_COLLECTION)
	count, err := col.CountDocuments(context.TODO(), bson.M{"token": token})
	return count > 0, err
}
